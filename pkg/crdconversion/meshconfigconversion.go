package crdconversion

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// serveMeshConfigConversion servers endpoint for the converter defined as convertMeshConfig function.
func serveMeshConfigConversion(w http.ResponseWriter, r *http.Request) {
	serve(w, r, convertMeshConfig)
}

// convertMeshConfig contains the business logic to convert meshconfigs.config.openservicemesh.io CRD
func convertMeshConfig(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	log.Info().Msgf("TEST converting meshconfig crd")

	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("conversion from a version to itself should not call the webhook: %s", toVersion)
	}

	log.Info().Msgf("TEST FromVersion meshconfig version %v", Object.GetAPIVersion())
	log.Info().Msgf("TEST ToVersion meshconfig version %v", toVersion)
	switch Object.GetAPIVersion() {
	case "config.openservicemesh.io/v1alpha1":
		switch toVersion {
		case "config.openservicemesh.io/v1alpha2":
			log.Info().Msgf("TEST converting v1alpha1 to v1alpha2")
			//log.Info().Msgf("object %v", convertedObject.Object)
			featureFlags, ok, _ := unstructured.NestedMap(convertedObject.Object, "spec", "featureFlags")
			log.Info().Msgf("TEST Featureflags : %v", featureFlags)
			if ok {
				delete(convertedObject.Object, "featureFlags")
				featureFlags["enableTestMode"] = "VALUE CHANGED BY CONVERTER"
				log.Info().Msgf("Featureflags : %v", featureFlags)
				if err := unstructured.SetNestedMap(convertedObject.Object, featureFlags, "spec", "featureFlags"); err != nil {
					log.Info().Msgf("TEST unable to set object")
					return nil, statusErrorWithMessage("unable to set object for version conversion version %q", toVersion)
				}
			}
		default:
			return nil, statusErrorWithMessage("unexpected conversion version %q", toVersion)
		}
	case "config.openservicemesh.io/v1alpha2":
		switch toVersion {
		case "config.openservicemesh.io/v1alpha1":
			log.Info().Msgf("TEST converting v1alpha2 to v1alpha1")
			//log.Info().Msgf("object %v", convertedObject.Object)
			featureFlags, ok, _ := unstructured.NestedMap(convertedObject.Object, "spec", "featureFlags")
			log.Info().Msgf("TEST Featureflags : %v", featureFlags)
			if ok {
				delete(convertedObject.Object, "featureFlags")
				delete(featureFlags, "enableTestMode")
				log.Info().Msgf("Featureflags : %v", featureFlags)
				if err := unstructured.SetNestedMap(convertedObject.Object, featureFlags, "spec", "featureFlags"); err != nil {
					log.Info().Msgf("TEST unable to set object")
					return nil, statusErrorWithMessage("unable to set object for version conversion version %q", toVersion)
				}
			}
		default:
			return nil, statusErrorWithMessage("unexpected conversion version %q", toVersion)
		}

	default:
		return nil, statusErrorWithMessage("unexpected conversion version %q", fromVersion)
	}
	log.Info().Msgf("TEST successfully Converted object %v", convertedObject.Object)
	return convertedObject, statusSucceed()
}
