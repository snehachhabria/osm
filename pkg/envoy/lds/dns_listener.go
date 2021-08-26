package lds

import (
	"sort"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
	xds_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	xds_listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	dnstable "github.com/envoyproxy/go-control-plane/envoy/data/dns/v3"
	dnsfilter "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/udp/dns_filter/v3alpha"
	stringmatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"

	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/errcode"
)

const resolverTimeout = 10 * time.Second

var knownSuffixes = []*stringmatcher.StringMatcher{
	{
		MatchPattern: &stringmatcher.StringMatcher_SafeRegex{
			SafeRegex: &stringmatcher.RegexMatcher{
				EngineType: &stringmatcher.RegexMatcher_GoogleRe2{GoogleRe2: &stringmatcher.RegexMatcher_GoogleRE2{}},
				Regex:      ".*", // Match everything.. All DNS queries go through Envoy. Unknown ones will be forwarded
			},
		},
	},
}

func (lb *listenerBuilder) newDNSListener() (*xds_listener.Listener, error) {
	address := envoy.GetAddress(constants.WildcardIPAddr, constants.EnvoyDNSListenerPort)
	// Convert the address to a UDP address
	address.GetSocketAddress().Protocol = xds_core.SocketAddress_UDP

	inlineDNSTable, err := lb.getInlineDNSTable()

	/*if len(inlineDNSTable.VirtualDomains) == 0 || err != nil {
		// No virtual domains or an error computing virtual domains
		return nil, nil
	}*/

	if err != nil {
		return nil, nil
	}

	dnsFilterConfig := &dnsfilter.DnsFilterConfig{
		StatPrefix: "dns",
		/*ServerConfig: &dnsfilter.DnsFilterConfig_ServerContextConfig{
			ConfigSource: &dnsfilter.DnsFilterConfig_ServerContextConfig_InlineDnsTable{InlineDnsTable: inlineDNSTable},
		},*/
		ClientConfig: &dnsfilter.DnsFilterConfig_ClientContextConfig{
			ResolverTimeout: ptypes.DurationProto(resolverTimeout),
			// We configure upstream resolver to resolver that always returns that it could not find the domain (NXDOMAIN)
			// As for this moment there is no setting to disable upstream resolving.
			UpstreamResolvers: []*xds_core.Address{{Address: &xds_core.Address_SocketAddress{
				SocketAddress: &xds_core.SocketAddress{
					//Protocol: xds_core.SocketAddress_UDP,
					Address: "10.0.0.10",
					PortSpecifier: &xds_core.SocketAddress_PortValue{
						PortValue: 53,
					},
				},
			}}},
			MaxPendingLookups: 256,
		},
	}

	if len(inlineDNSTable.VirtualDomains) != 0 {
		dnsFilterConfig.ServerConfig =
			&dnsfilter.DnsFilterConfig_ServerContextConfig{
				ConfigSource: &dnsfilter.DnsFilterConfig_ServerContextConfig_InlineDnsTable{InlineDnsTable: inlineDNSTable},
			}
	}

	dnsFilterConfigMarshal, err := ptypes.MarshalAny(dnsFilterConfig)
	if err != nil {
		log.Error().Err(err).Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrMarshallingXDSResource)).
			Msgf("TEST Error marshalling HttpConnectionManager object")
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

func (lb *listenerBuilder) getInlineDNSTable() (*dnstable.DnsTable, error) {
	upstreamServices := lb.meshCatalog.ListMeshServicesForIdentity(lb.serviceIdentity)
	// build a virtual domain for each service visible to this sidecar
	virtualDomains := make([]*dnstable.DnsTable_DnsVirtualDomain, 0)
	log.Info().Msgf("TEST build DNS Table for proxy %s with services %v", lb.serviceIdentity, upstreamServices)

	if len(upstreamServices) == 0 {
		log.Debug().Msgf("Proxy with identity %s does not have any allowed upstream services", lb.serviceIdentity)
		return &dnstable.DnsTable{
			VirtualDomains: virtualDomains,
			KnownSuffixes:  knownSuffixes,
		}, nil
	}

	for _, upstreamSvc := range upstreamServices {
		log.Trace().Msgf("TEST Building dns filter chain for upstream service %s for proxy with identity %s", upstreamSvc, lb.serviceIdentity)
		protocolToPortMap, err := lb.meshCatalog.GetPortToProtocolMappingForService(upstreamSvc)
		if err != nil {
			log.Error().Err(err).Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrGettingServicePorts)).
				Msgf("Error retrieving port to protocol mapping for upstream service %s", upstreamSvc)
			continue
		}

		for port, appProtocol := range protocolToPortMap {
			switch strings.ToLower(appProtocol) {
			case constants.ProtocolHTTP, constants.ProtocolGRPC:
				endpoints, err := lb.meshCatalog.GetResolvableServiceEndpoints(upstreamSvc)
				if err != nil {
					log.Error().Err(err).Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrGettingResolvableServiceEndpoints)).
						Msgf("Error getting GetResolvableServiceEndpoints for %q", upstreamSvc)
					return &dnstable.DnsTable{}, err
				}

				if len(endpoints) == 0 {
					err := errors.Errorf("Endpoints not found for service %q", upstreamSvc)
					log.Error().Err(err).Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrEndpointsNotFound)).
						Msgf("Error constructing HTTP filter chain match for service %q", upstreamSvc)
					return nil, err
				}

				endpointSet := mapset.NewSet()
				for _, endp := range endpoints {
					endpointSet.Add(endp.IP.String())
				}

				// For deterministic ordering
				var sortedEndpoints []string
				endpointSet.Each(func(elem interface{}) bool {
					sortedEndpoints = append(sortedEndpoints, elem.(string))
					return false
				})
				sort.Strings(sortedEndpoints)

				/*virtualDomains = append(virtualDomains, &dnstable.DnsTable_DnsVirtualDomain{
					Name: upstreamSvc.Name + "." + upstreamSvc.Namespace + ".svc.cluster.local",
					Endpoint: &dnstable.DnsTable_DnsEndpoint{
						EndpointConfig: &dnstable.DnsTable_DnsEndpoint_AddressList{
							AddressList: &dnstable.DnsTable_AddressList{Address: sortedEndpoints},
						},
					},
				})*/

				// Add endpoints to gateway's in other clusters
				endpoints, err = lb.meshCatalog.GetMulticlusterGatewayEndpoints(upstreamSvc)
				if err != nil {
					log.Error().Err(err).Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrGettingMulticlusterGatewayEndpoints)).
						Msgf("Error getting GetMulticlusterGatewaysEndpoints for %q", upstreamSvc)
					return &dnstable.DnsTable{}, err
				}

				if len(endpoints) == 0 {
					err := errors.Errorf("Multicluster endpoints not found for service %q", upstreamSvc)
					log.Error().Err(err).Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrEndpointsNotFound)).
						Msgf("Error constructing HTTP filter chain match for service %q", upstreamSvc)
					return nil, err
				}
				//var sortedEndpoints []string
				endpointSet = mapset.NewSet()
				for _, endp := range endpoints {
					endpointSet.Add(endp.IP.String())
				}

				// For deterministic ordering
				endpointSet.Each(func(elem interface{}) bool {
					sortedEndpoints = append(sortedEndpoints, elem.(string))
					return false
				})
				sort.Strings(sortedEndpoints)

				virtualDomains = append(virtualDomains, &dnstable.DnsTable_DnsVirtualDomain{
					Name: upstreamSvc.Name + "." + "osmmesh",
					Endpoint: &dnstable.DnsTable_DnsEndpoint{
						EndpointConfig: &dnstable.DnsTable_DnsEndpoint_AddressList{
							AddressList: &dnstable.DnsTable_AddressList{Address: sortedEndpoints},
						},
					},
				})

			case constants.ProtocolTCP:
				continue

			default:
				log.Error().Msgf("Cannot build address list, unsupported protocol %s for upstream:port %s:%d", appProtocol, upstreamSvc, port)
			}
		}

	}

	sort.Stable(DnsTableByName(virtualDomains)) // for stable Envoy config

	dnsTable := &dnstable.DnsTable{
		VirtualDomains:     virtualDomains,
		ExternalRetryCount: 0,
		KnownSuffixes:      knownSuffixes,
	}

	log.Info().Msgf("TEST DNS Table for proxy %s is %v", lb.serviceIdentity, dnsTable)

	return dnsTable, nil
}

type DnsTableByName []*dnstable.DnsTable_DnsVirtualDomain

func (a DnsTableByName) Len() int      { return len(a) }
func (a DnsTableByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a DnsTableByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}
