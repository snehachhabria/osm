// Package injector implements OSM's automatic sidecar injection facility. The sidecar injector's mutating webhook
// admission controller intercepts pod creation requests to mutate the pod spec to inject the sidecar proxy.
package crdconverter

import (
	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/logger"
)

const (
	envoyBootstrapConfigVolume = "envoy-bootstrap-config-volume"
)

var log = logger.New("crd-converter")

// converterWebhook is the type used to represent the webhook for crd conversion
type converterWebhook struct {
	config Config
	cert   certificate.Certificater
}

// Config is the type used to represent the config options for the sidecar injection
type Config struct {
	// ListenPort defines the port on which the sidecar injector listens
	ListenPort int
}
