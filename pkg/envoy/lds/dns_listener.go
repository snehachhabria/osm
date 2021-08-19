package lds

import (
	"time"

	xds_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	xds_listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	dnstable "github.com/envoyproxy/go-control-plane/envoy/data/dns/v3"
	dnsfilter "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/udp/dns_filter/v3alpha"
	"github.com/golang/protobuf/ptypes"

	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/errcode"
)

const resolverTimeout = 10 * time.Second

func newDNSListener() (*xds_listener.Listener, error) {
	address := envoy.GetAddress(constants.WildcardIPAddr, constants.EnvoyDNSListenerPort)
	// Convert the address to a UDP address
	address.GetSocketAddress().Protocol = xds_core.SocketAddress_UDP

	inlineDNSTable := buildInlineDNSTable(node, push)

	dnsFilterConfig := &dnsfilter.DnsFilterConfig{
		StatPrefix: "dns",
		ServerConfig: &dnsfilter.DnsFilterConfig_ServerContextConfig{
			ConfigSource: &dnsfilter.DnsFilterConfig_ServerContextConfig_InlineDnsTable{InlineDnsTable: inlineDNSTable},
		},
		ClientConfig: &dnsfilter.DnsFilterConfig_ClientContextConfig{
			ResolverTimeout: ptypes.DurationProto(resolverTimeout),
			// no upstream resolves. Envoy will use the ambient ones
			MaxPendingLookups: 256, // arbitrary
		},
	}
	dnsFilterConfigMarshal, err := ptypes.MarshalAny(dnsFilterConfig)
	if err != nil {
		log.Error().Err(err).Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrMarshallingXDSResource)).
			Msgf("Error marshalling HttpConnectionManager object")
		return nil, err
	}

	dnsFilter := &xds_listener.ListenerFilter{
		Name: "envoy.filters.udp.dns_filter",
		ConfigType: &xds_listener.ListenerFilter_TypedConfig{
			TypedConfig: dnsFilterConfigMarshal,
		},
	}

	return &xds_listener.Listener{
		Name:             dnsListenerName,
		Address:          address,
		ListenerFilters:  []*xds_listener.ListenerFilter{dnsFilter},
		TrafficDirection: xds_core.TrafficDirection_OUTBOUND,
		// DNS listener requires SO_REUSEPORT option to be set esp when concurrency >1
		ReusePort: true,
	}, nil
}

func buildInlineDNSTable(node *model.Proxy, push *model.PushContext) *dnstable.DnsTable {

	// build a virtual domain for each service visible to this sidecar
	virtualDomains := make([]*dnstable.DnsTable_DnsVirtualDomain, 0)

	for _, svc := range push.Services(node) {
		// we cannot take services with wildcards in the address field. The reason
		// is that even if we provide some dummy IP (subject to enabling this
		// feature in Envoy), after capturing the traffic from the app, the
		// sidecar would need to forward to the real IP. But to determine the real
		// IP, the sidecar would have to know the non-wildcard FQDN that the
		// application was trying to resolve. This information is not available
		// for TCP services. The wildcard hostname is not a problem for HTTP
		// services though, as we usually setup a listener on 0.0.0.0, process
		// based on http virtual host and forward to the orig destination IP.
		//
		// Long story short, if the user has a TCP service of the form
		//
		// host: *.mysql.aws.com, port 3306,
		//
		// our only recourse is to allocate a 0.0.0.0:3306 passthrough listener and forward to
		// original dest IP. It is now the user's responsibility to not allocate
		// another wildcard service on the same port. i.e.
		//
		// 1. host: *.mysql.aws.com, port 3306
		// 2. host: *.mongo.aws.com, port 3306 will result in conflict.
		//
		// Traffic will still flow but metrics wont be correct
		// as two different TCP services are consuming the
		// same wildcard passthrough TCP listener 0.0.0.0:3306.
		//
		if svc.Hostname.IsWildCarded() {
			continue
		}

		svcAddress := svc.GetServiceAddressForProxy(node)
		var addressList []string

		// The IP will be unspecified here if its headless service or if the auto
		// IP allocation logic for service entry was unable to allocate an IP.
		if svcAddress == constants.UnspecifiedIP {
			// For all k8s headless services, populate the dns table with the endpoint IPs as k8s does.
			// TODO: Need to have an entry per pod hostname of stateful set but for this, we need to parse
			// the stateful set object, associate the object with the appropriate kubernetes headless service
			// and then derive the stable network identities.
			if svc.Attributes.ServiceRegistry == string(serviceregistry.Kubernetes) &&
				svc.Resolution == model.Passthrough && len(svc.Ports) > 0 {
				// TODO: this is used in two places now. Needs to be cached as part of the headless service
				// object to avoid the costly lookup in the registry code
				if instances, err := push.InstancesByPort(svc, svc.Ports[0].Port, nil); err == nil {
					for _, instance := range instances {
						// TODO: should we skip the node's own IP like we do in listener?
						addressList = append(addressList, instance.Endpoint.Address)
					}
				}
			}

			if len(addressList) == 0 {
				// could not reliably determine the addresses of endpoints of headless service
				// or this is not a k8s service
				continue
			}
		} else {
			addressList = append(addressList, svcAddress)
		}

		virtualDomains = append(virtualDomains, &dnstable.DnsTable_DnsVirtualDomain{
			Name: string(svc.Hostname),
			Endpoint: &dnstable.DnsTable_DnsEndpoint{
				EndpointConfig: &dnstable.DnsTable_DnsEndpoint_AddressList{
					AddressList: &dnstable.DnsTable_AddressList{Address: addressList},
				},
			},
		})

		// If this is a kubernetes service, generate short form names (name.namespace) and
		// just name (if proxy is in same namespace).
		if svc.Attributes.ServiceRegistry == string(serviceregistry.Kubernetes) {
			virtualDomains = append(virtualDomains, &dnstable.DnsTable_DnsVirtualDomain{
				Name: svc.Attributes.Name + "." + svc.Attributes.Namespace,
				Endpoint: &dnstable.DnsTable_DnsEndpoint{
					EndpointConfig: &dnstable.DnsTable_DnsEndpoint_AddressList{
						AddressList: &dnstable.DnsTable_AddressList{Address: addressList},
					},
				},
			})
			if node.ConfigNamespace == svc.Attributes.Namespace {
				virtualDomains = append(virtualDomains, &dnstable.DnsTable_DnsVirtualDomain{
					Name: svc.Attributes.Name,
					Endpoint: &dnstable.DnsTable_DnsEndpoint{
						EndpointConfig: &dnstable.DnsTable_DnsEndpoint_AddressList{
							AddressList: &dnstable.DnsTable_AddressList{Address: addressList},
						},
					},
				})
			}
		}
	}

	return &dnstable.DnsTable{
		VirtualDomains: virtualDomains,
		KnownSuffixes:  knownSuffixes,
	}
}
