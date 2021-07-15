package crdconversion

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// serveTCPRouteConversion servers endpoint for the converter defined as convertTCPRoute function.
func serveTCPRouteConversion(w http.ResponseWriter, r *http.Request) {
	serve(w, r, convertTCPRoute)
}

// convertTCPRoute contains the business logic to convert tcproutes.specs.smi-spec.io CRD
// Example implementation reference : https://github.com/kubernetes/kubernetes/blob/release-1.21/test/images/agnhost/crd-conversion-webhook/converter/example_converter.go
func convertTCPRoute(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	//log.Info().Msgf("TEST converting TCPRoute crd")

	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("conversion from a version to itself should not call the webhook: %s", toVersion)
	}
	//log.Info().Msgf("TEST successfully converted TCPRoute object")
	return convertedObject, statusSucceed()
}
