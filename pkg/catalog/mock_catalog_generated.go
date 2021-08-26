// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openservicemesh/osm/pkg/catalog (interfaces: MeshCataloger)

// Package catalog is a generated GoMock package.
package catalog

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	endpoint "github.com/openservicemesh/osm/pkg/endpoint"
	identity "github.com/openservicemesh/osm/pkg/identity"
	k8s "github.com/openservicemesh/osm/pkg/k8s"
	service "github.com/openservicemesh/osm/pkg/service"
	trafficpolicy "github.com/openservicemesh/osm/pkg/trafficpolicy"
)

// MockMeshCataloger is a mock of MeshCataloger interface
type MockMeshCataloger struct {
	ctrl     *gomock.Controller
	recorder *MockMeshCatalogerMockRecorder
}

// MockMeshCatalogerMockRecorder is the mock recorder for MockMeshCataloger
type MockMeshCatalogerMockRecorder struct {
	mock *MockMeshCataloger
}

// NewMockMeshCataloger creates a new mock instance
func NewMockMeshCataloger(ctrl *gomock.Controller) *MockMeshCataloger {
	mock := &MockMeshCataloger{ctrl: ctrl}
	mock.recorder = &MockMeshCatalogerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMeshCataloger) EXPECT() *MockMeshCatalogerMockRecorder {
	return m.recorder
}

// GetEgressTrafficPolicy mocks base method
func (m *MockMeshCataloger) GetEgressTrafficPolicy(arg0 identity.ServiceIdentity) (*trafficpolicy.EgressTrafficPolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEgressTrafficPolicy", arg0)
	ret0, _ := ret[0].(*trafficpolicy.EgressTrafficPolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEgressTrafficPolicy indicates an expected call of GetEgressTrafficPolicy
func (mr *MockMeshCatalogerMockRecorder) GetEgressTrafficPolicy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEgressTrafficPolicy", reflect.TypeOf((*MockMeshCataloger)(nil).GetEgressTrafficPolicy), arg0)
}

// GetIngressTrafficPolicy mocks base method
func (m *MockMeshCataloger) GetIngressTrafficPolicy(arg0 service.MeshService) (*trafficpolicy.IngressTrafficPolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIngressTrafficPolicy", arg0)
	ret0, _ := ret[0].(*trafficpolicy.IngressTrafficPolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIngressTrafficPolicy indicates an expected call of GetIngressTrafficPolicy
func (mr *MockMeshCatalogerMockRecorder) GetIngressTrafficPolicy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIngressTrafficPolicy", reflect.TypeOf((*MockMeshCataloger)(nil).GetIngressTrafficPolicy), arg0)
}

// GetKubeController mocks base method
func (m *MockMeshCataloger) GetKubeController() k8s.Controller {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKubeController")
	ret0, _ := ret[0].(k8s.Controller)
	return ret0
}

// GetKubeController indicates an expected call of GetKubeController
func (mr *MockMeshCatalogerMockRecorder) GetKubeController() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKubeController", reflect.TypeOf((*MockMeshCataloger)(nil).GetKubeController))
}

// GetMulticlusterGatewayEndpoints mocks base method
func (m *MockMeshCataloger) GetMulticlusterGatewayEndpoints(arg0 service.MeshService) ([]endpoint.Endpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMulticlusterGatewayEndpoints", arg0)
	ret0, _ := ret[0].([]endpoint.Endpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMulticlusterGatewayEndpoints indicates an expected call of GetMulticlusterGatewayEndpoints
func (mr *MockMeshCatalogerMockRecorder) GetMulticlusterGatewayEndpoints(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMulticlusterGatewayEndpoints", reflect.TypeOf((*MockMeshCataloger)(nil).GetMulticlusterGatewayEndpoints), arg0)
}

// GetPortToProtocolMappingForService mocks base method
func (m *MockMeshCataloger) GetPortToProtocolMappingForService(arg0 service.MeshService) (map[uint32]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPortToProtocolMappingForService", arg0)
	ret0, _ := ret[0].(map[uint32]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPortToProtocolMappingForService indicates an expected call of GetPortToProtocolMappingForService
func (mr *MockMeshCatalogerMockRecorder) GetPortToProtocolMappingForService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPortToProtocolMappingForService", reflect.TypeOf((*MockMeshCataloger)(nil).GetPortToProtocolMappingForService), arg0)
}

// GetResolvableServiceEndpoints mocks base method
func (m *MockMeshCataloger) GetResolvableServiceEndpoints(arg0 service.MeshService) ([]endpoint.Endpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResolvableServiceEndpoints", arg0)
	ret0, _ := ret[0].([]endpoint.Endpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResolvableServiceEndpoints indicates an expected call of GetResolvableServiceEndpoints
func (mr *MockMeshCatalogerMockRecorder) GetResolvableServiceEndpoints(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResolvableServiceEndpoints", reflect.TypeOf((*MockMeshCataloger)(nil).GetResolvableServiceEndpoints), arg0)
}

// GetServiceHostnames mocks base method
func (m *MockMeshCataloger) GetServiceHostnames(arg0 service.MeshService, arg1 service.Locality) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceHostnames", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceHostnames indicates an expected call of GetServiceHostnames
func (mr *MockMeshCatalogerMockRecorder) GetServiceHostnames(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceHostnames", reflect.TypeOf((*MockMeshCataloger)(nil).GetServiceHostnames), arg0, arg1)
}

// GetTargetPortToProtocolMappingForService mocks base method
func (m *MockMeshCataloger) GetTargetPortToProtocolMappingForService(arg0 service.MeshService) (map[uint32]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTargetPortToProtocolMappingForService", arg0)
	ret0, _ := ret[0].(map[uint32]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTargetPortToProtocolMappingForService indicates an expected call of GetTargetPortToProtocolMappingForService
func (mr *MockMeshCatalogerMockRecorder) GetTargetPortToProtocolMappingForService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTargetPortToProtocolMappingForService", reflect.TypeOf((*MockMeshCataloger)(nil).GetTargetPortToProtocolMappingForService), arg0)
}

// GetWeightedClustersForUpstream mocks base method
func (m *MockMeshCataloger) GetWeightedClustersForUpstream(arg0 service.MeshService) []service.WeightedCluster {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWeightedClustersForUpstream", arg0)
	ret0, _ := ret[0].([]service.WeightedCluster)
	return ret0
}

// GetWeightedClustersForUpstream indicates an expected call of GetWeightedClustersForUpstream
func (mr *MockMeshCatalogerMockRecorder) GetWeightedClustersForUpstream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWeightedClustersForUpstream", reflect.TypeOf((*MockMeshCataloger)(nil).GetWeightedClustersForUpstream), arg0)
}

// ListEndpointsForServiceIdentity mocks base method
func (m *MockMeshCataloger) ListEndpointsForServiceIdentity(arg0 identity.ServiceIdentity, arg1 service.MeshService) ([]endpoint.Endpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEndpointsForServiceIdentity", arg0, arg1)
	ret0, _ := ret[0].([]endpoint.Endpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEndpointsForServiceIdentity indicates an expected call of ListEndpointsForServiceIdentity
func (mr *MockMeshCatalogerMockRecorder) ListEndpointsForServiceIdentity(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEndpointsForServiceIdentity", reflect.TypeOf((*MockMeshCataloger)(nil).ListEndpointsForServiceIdentity), arg0, arg1)
}

// ListInboundServiceIdentities mocks base method
func (m *MockMeshCataloger) ListInboundServiceIdentities(arg0 identity.ServiceIdentity) ([]identity.ServiceIdentity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListInboundServiceIdentities", arg0)
	ret0, _ := ret[0].([]identity.ServiceIdentity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListInboundServiceIdentities indicates an expected call of ListInboundServiceIdentities
func (mr *MockMeshCatalogerMockRecorder) ListInboundServiceIdentities(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListInboundServiceIdentities", reflect.TypeOf((*MockMeshCataloger)(nil).ListInboundServiceIdentities), arg0)
}

// ListInboundTrafficPolicies mocks base method
func (m *MockMeshCataloger) ListInboundTrafficPolicies(arg0 identity.ServiceIdentity, arg1 []service.MeshService) []*trafficpolicy.InboundTrafficPolicy {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListInboundTrafficPolicies", arg0, arg1)
	ret0, _ := ret[0].([]*trafficpolicy.InboundTrafficPolicy)
	return ret0
}

// ListInboundTrafficPolicies indicates an expected call of ListInboundTrafficPolicies
func (mr *MockMeshCatalogerMockRecorder) ListInboundTrafficPolicies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListInboundTrafficPolicies", reflect.TypeOf((*MockMeshCataloger)(nil).ListInboundTrafficPolicies), arg0, arg1)
}

// ListInboundTrafficTargetsWithRoutes mocks base method
func (m *MockMeshCataloger) ListInboundTrafficTargetsWithRoutes(arg0 identity.ServiceIdentity) ([]trafficpolicy.TrafficTargetWithRoutes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListInboundTrafficTargetsWithRoutes", arg0)
	ret0, _ := ret[0].([]trafficpolicy.TrafficTargetWithRoutes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListInboundTrafficTargetsWithRoutes indicates an expected call of ListInboundTrafficTargetsWithRoutes
func (mr *MockMeshCatalogerMockRecorder) ListInboundTrafficTargetsWithRoutes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListInboundTrafficTargetsWithRoutes", reflect.TypeOf((*MockMeshCataloger)(nil).ListInboundTrafficTargetsWithRoutes), arg0)
}

// ListMeshServicesForIdentity mocks base method
func (m *MockMeshCataloger) ListMeshServicesForIdentity(arg0 identity.ServiceIdentity) []service.MeshService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMeshServicesForIdentity", arg0)
	ret0, _ := ret[0].([]service.MeshService)
	return ret0
}

// ListMeshServicesForIdentity indicates an expected call of ListMeshServicesForIdentity
func (mr *MockMeshCatalogerMockRecorder) ListMeshServicesForIdentity(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMeshServicesForIdentity", reflect.TypeOf((*MockMeshCataloger)(nil).ListMeshServicesForIdentity), arg0)
}

// ListOutboundServiceIdentities mocks base method
func (m *MockMeshCataloger) ListOutboundServiceIdentities(arg0 identity.ServiceIdentity) ([]identity.ServiceIdentity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOutboundServiceIdentities", arg0)
	ret0, _ := ret[0].([]identity.ServiceIdentity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOutboundServiceIdentities indicates an expected call of ListOutboundServiceIdentities
func (mr *MockMeshCatalogerMockRecorder) ListOutboundServiceIdentities(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOutboundServiceIdentities", reflect.TypeOf((*MockMeshCataloger)(nil).ListOutboundServiceIdentities), arg0)
}

// ListOutboundServicesForIdentity mocks base method
func (m *MockMeshCataloger) ListOutboundServicesForIdentity(arg0 identity.ServiceIdentity) []service.MeshService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOutboundServicesForIdentity", arg0)
	ret0, _ := ret[0].([]service.MeshService)
	return ret0
}

// ListOutboundServicesForIdentity indicates an expected call of ListOutboundServicesForIdentity
func (mr *MockMeshCatalogerMockRecorder) ListOutboundServicesForIdentity(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOutboundServicesForIdentity", reflect.TypeOf((*MockMeshCataloger)(nil).ListOutboundServicesForIdentity), arg0)
}

// ListOutboundServicesForMulticlusterGateway mocks base method
func (m *MockMeshCataloger) ListOutboundServicesForMulticlusterGateway() []service.MeshService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOutboundServicesForMulticlusterGateway")
	ret0, _ := ret[0].([]service.MeshService)
	return ret0
}

// ListOutboundServicesForMulticlusterGateway indicates an expected call of ListOutboundServicesForMulticlusterGateway
func (mr *MockMeshCatalogerMockRecorder) ListOutboundServicesForMulticlusterGateway() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOutboundServicesForMulticlusterGateway", reflect.TypeOf((*MockMeshCataloger)(nil).ListOutboundServicesForMulticlusterGateway))
}

// ListOutboundTrafficPolicies mocks base method
func (m *MockMeshCataloger) ListOutboundTrafficPolicies(arg0 identity.ServiceIdentity) []*trafficpolicy.OutboundTrafficPolicy {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOutboundTrafficPolicies", arg0)
	ret0, _ := ret[0].([]*trafficpolicy.OutboundTrafficPolicy)
	return ret0
}

// ListOutboundTrafficPolicies indicates an expected call of ListOutboundTrafficPolicies
func (mr *MockMeshCatalogerMockRecorder) ListOutboundTrafficPolicies(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOutboundTrafficPolicies", reflect.TypeOf((*MockMeshCataloger)(nil).ListOutboundTrafficPolicies), arg0)
}

// ListServiceIdentitiesForService mocks base method
func (m *MockMeshCataloger) ListServiceIdentitiesForService(arg0 service.MeshService) ([]identity.ServiceIdentity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListServiceIdentitiesForService", arg0)
	ret0, _ := ret[0].([]identity.ServiceIdentity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListServiceIdentitiesForService indicates an expected call of ListServiceIdentitiesForService
func (mr *MockMeshCatalogerMockRecorder) ListServiceIdentitiesForService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServiceIdentitiesForService", reflect.TypeOf((*MockMeshCataloger)(nil).ListServiceIdentitiesForService), arg0)
}
