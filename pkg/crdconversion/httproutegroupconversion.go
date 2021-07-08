package crdconversion

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// serveHTTPRouteGroupConversion servers endpoint for the converter defined as convertHTTPRouteGroup function.
func serveHTTPRouteGroupConversion(w http.ResponseWriter, r *http.Request) {
	serve(w, r, convertHTTPRouteGroup)
}

// convertEgressPolicy contains the business logic to convert httproutegroups.specs.smi-spec.io CRD
// Example implementation reference : https://github.com/kubernetes/kubernetes/blob/release-1.21/test/images/agnhost/crd-conversion-webhook/converter/example_converter.go
func convertHTTPRouteGroup(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	//log.Info().Msgf("TEST converting HTTPRouteGroup crd")

	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	//log.Info().Msgf("TEST HTTPRouteGroup FROMversion %s TOversion %s object %v", fromVersion, toVersion, convertedObject.Object)

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("conversion from a version to itself should not call the webhook: %s", toVersion)
	}
	//log.Info().Msgf("TEST successfully converted HTTPRouteGroup object")
	return convertedObject, statusSucceed()
}
