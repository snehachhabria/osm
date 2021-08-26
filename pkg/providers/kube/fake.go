package kube

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/openservicemesh/osm/pkg/endpoint"
	"github.com/openservicemesh/osm/pkg/identity"
	"github.com/openservicemesh/osm/pkg/service"
	"github.com/openservicemesh/osm/pkg/tests"
)

// Provider interface combines endpoint.Provider and service.Provider
type Provider interface {
	endpoint.Provider
	service.Provider
}

// NewFakeProvider implements mesh.EndpointsProvider, which creates a test Kubernetes cluster/compute provider.
func NewFakeProvider() Provider {
	return fakeClient{
		endpoints: map[string][]endpoint.Endpoint{
			tests.BookstoreV1Service.String():   {tests.Endpoint},
			tests.BookstoreV2Service.String():   {tests.Endpoint},
			tests.BookbuyerService.String():     {tests.Endpoint},
			tests.BookstoreApexService.String(): {tests.Endpoint},
		},
		gatewayEndpoints: map[string][]endpoint.Endpoint{
			tests.BookstoreV1Service.String():   {tests.GatewayEndpoint},
			tests.BookstoreV2Service.String():   {tests.GatewayEndpoint},
			tests.BookbuyerService.String():     {tests.GatewayEndpoint},
			tests.BookstoreApexService.String(): {tests.GatewayEndpoint},
		},
		services: map[identity.K8sServiceAccount][]service.MeshService{
			tests.BookstoreServiceAccount:   {tests.BookstoreV1Service, tests.BookstoreApexService},
			tests.BookstoreV2ServiceAccount: {tests.BookstoreV2Service},
			tests.BookbuyerServiceAccount:   {tests.BookbuyerService},
		},
		svcAccountEndpoints: map[identity.K8sServiceAccount][]endpoint.Endpoint{
			tests.BookstoreServiceAccount:   {tests.Endpoint, tests.Endpoint},
			tests.BookstoreV2ServiceAccount: {tests.Endpoint},
			tests.BookbuyerServiceAccount:   {tests.Endpoint},
		},
	}
}

type fakeClient struct {
	endpoints           map[string][]endpoint.Endpoint
	gatewayEndpoints    map[string][]endpoint.Endpoint
	services            map[identity.K8sServiceAccount][]service.MeshService
	svcAccountEndpoints map[identity.K8sServiceAccount][]endpoint.Endpoint
}

// ListEndpointsForService retrieves the IP addresses comprising the given service.
func (f fakeClient) ListEndpointsForService(svc service.MeshService) []endpoint.Endpoint {
	if svc, ok := f.endpoints[svc.String()]; ok {
		return svc
	}
	panic(fmt.Sprintf("You are asking for MeshService=%s but the fake Kubernetes client has not been initialized with this. What we have is: %+v", svc, f.endpoints))
}

// ListEndpointsForIdentity retrieves the IP addresses comprising the given service account.
// Note: ServiceIdentity must be in the format "name.namespace" [https://github.com/openservicemesh/osm/issues/3188]
func (f fakeClient) ListEndpointsForIdentity(serviceIdentity identity.ServiceIdentity) []endpoint.Endpoint {
	sa := serviceIdentity.ToK8sServiceAccount()
	if ep, ok := f.svcAccountEndpoints[sa]; ok {
		return ep
	}
	panic(fmt.Sprintf("You are asking for K8sServiceAccount=%s but the fake Kubernetes client has not been initialized with this. What we have is: %+v", sa, f.svcAccountEndpoints))
}

func (f fakeClient) GetServicesForServiceIdentity(serviceIdentity identity.ServiceIdentity) ([]service.MeshService, error) {
	sa := serviceIdentity.ToK8sServiceAccount()
	services, ok := f.services[sa]
	if !ok {
		return nil, errors.Errorf("ServiceAccount %s is not in cache: %+v", sa, f.services)
	}
	return services, nil
}

func (f fakeClient) ListServices() ([]service.MeshService, error) {
	var services []service.MeshService

	for _, svcs := range f.services {
		services = append(services, svcs...)
	}
	return services, nil
}

func (f fakeClient) ListServiceIdentitiesForService(svc service.MeshService) ([]identity.ServiceIdentity, error) {
	var serviceIdentities []identity.ServiceIdentity

	for svcID := range f.services {
		serviceIdentities = append(serviceIdentities, svcID.ToServiceIdentity())
	}
	return serviceIdentities, nil
}

func (f fakeClient) GetHostnamesForService(svc service.MeshService, locality service.Locality) ([]string, error) {
	var hostnames []string

	serviceName := svc.Name
	namespace := svc.Namespace
	port := tests.Endpoint.Port

	if locality == service.LocalNS {
		hostnames = append(hostnames, serviceName)
	}

	hostnames = append(hostnames, fmt.Sprintf("%s.%s", serviceName, namespace))                   // service.namespace
	hostnames = append(hostnames, fmt.Sprintf("%s.%s.svc", serviceName, namespace))               // service.namespace.svc
	hostnames = append(hostnames, fmt.Sprintf("%s.%s.svc.cluster", serviceName, namespace))       // service.namespace.svc.cluster
	hostnames = append(hostnames, fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)) // service.namespace.svc.cluster.local
	if locality == service.LocalNS {
		// Within the same namespace, service name is resolvable to its address
		hostnames = append(hostnames, fmt.Sprintf("%s:%d", serviceName, port)) // service:port
	}
	hostnames = append(hostnames, fmt.Sprintf("%s.%s:%d", serviceName, namespace, port))                   // service.namespace:port
	hostnames = append(hostnames, fmt.Sprintf("%s.%s.svc:%d", serviceName, namespace, port))               // service.namespace.svc:port
	hostnames = append(hostnames, fmt.Sprintf("%s.%s.svc.cluster:%d", serviceName, namespace, port))       // service.namespace.svc.cluster:port
	hostnames = append(hostnames, fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)) // service.namespace.svc.cluster.local:port
	return hostnames, nil
}

func (f fakeClient) GetTargetPortToProtocolMappingForService(svc service.MeshService) (map[uint32]string, error) {
	return map[uint32]string{uint32(tests.Endpoint.Port): "http"}, nil
}

func (f fakeClient) GetPortToProtocolMappingForService(svc service.MeshService) (map[uint32]string, error) {
	return map[uint32]string{uint32(tests.Endpoint.Port): "http"}, nil
}

// GetID returns the unique identifier of the Provider.
func (f fakeClient) GetID() string {
	return "Fake Kubernetes Client"
}

func (f fakeClient) GetResolvableEndpointsForService(svc service.MeshService) ([]endpoint.Endpoint, error) {
	endpoints, found := f.endpoints[svc.String()]
	if !found {
		return nil, errServiceNotFound
	}
	return endpoints, nil
}

func (f fakeClient) GetMulticlusterEndpointsForService(svc service.MeshService) ([]endpoint.Endpoint, error) {
	endpoints, found := f.gatewayEndpoints[svc.String()]
	if !found {
		return nil, errServiceNotFound
	}
	return endpoints, nil
}
