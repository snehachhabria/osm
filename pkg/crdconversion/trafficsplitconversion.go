package crdconversion

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// serveTCPRoutesConversion servers endpoint for the converter defined as convertTrafficSplit function.
func serveTrafficSplitConversion(w http.ResponseWriter, r *http.Request) {
	serve(w, r, convertTrafficSplit)
}

// convertTrafficAccess contains the business logic to convert trafficsplits.access.smi-spec.io CRD
// Example implementation reference : https://github.com/kubernetes/kubernetes/blob/release-1.21/test/images/agnhost/crd-conversion-webhook/converter/example_converter.go
func convertTrafficSplit(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	//log.Info().Msgf("TEST converting TrafficSplit crd")

	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	//log.Info().Msgf("TEST TrafficSplit FROMversion %s TOversion %s object %v", fromVersion, toVersion, convertedObject.Object)

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("conversion from a version to itself should not call the webhook: %s", toVersion)
	}
	//log.Info().Msgf("TEST successfully converted TrafficSplit object")
	return convertedObject, statusSucceed()
}
