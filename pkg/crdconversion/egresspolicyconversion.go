package crdconversion

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// serveEgressPolicyConversion servers endpoint for the converter defined as convertEgressPolicy function.
func serveEgressPolicyConversion(w http.ResponseWriter, r *http.Request) {
	serve(w, r, convertEgressPolicy)
}

// convertEgressPolicy contains the business logic to convert egresses.policy.openservicemesh.io CRD
// Example implementation reference : https://github.com/kubernetes/kubernetes/blob/release-1.21/test/images/agnhost/crd-conversion-webhook/converter/example_converter.go
func convertEgressPolicy(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	//log.Info().Msgf("TEST converting egress policy crd")

	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("conversion from a version to itself should not call the webhook: %s", toVersion)
	}
	//log.Info().Msgf("TEST successfully converted EgressPolicy object")
	return convertedObject, statusSucceed()
}
