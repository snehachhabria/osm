package crdconverter

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/certificate/providers"
	"github.com/openservicemesh/osm/pkg/constants"
)

const (

	// webhookCreatePod is the HTTP path at which the webhook expects to receive mesh config crd conversion events
	webhookMeshConfigConvert = "/meshconfigconvert"

	// WebhookHealthPath is the HTTP path at which the health of the webhook can be queried
	WebhookHealthPath = "/healthz"

	// crdConvertedServiceName is the name of the OSM crd conversion webhook service
	crdConvertedServiceName = "osm-crd-conversion-webhook"
)

// NewConversionWebhook starts a new web server handling requests from the injector MutatingWebhookConfiguration
func NewConversionWebhook(config Config, kubeClient kubernetes.Interface, crdClient apiclient.ApiextensionsV1Interface, certManager certificate.Manager, osmNamespace string, stop <-chan struct{}) error {
	// This is a certificate issued for the crd webhook handler
	// This cert does not have to be related to the Envoy certs, but it does have to match
	// the cert provisioned with the MutatingWebhookConfiguration
	crdConverterwebhookHandlerCert, err := certManager.IssueCertificate(
		certificate.CommonName(fmt.Sprintf("%s.%s.svc", crdConvertedServiceName, osmNamespace)),
		constants.XDSCertificateValidityPeriod)
	if err != nil {
		return errors.Errorf("Error issuing certificate for the mutating webhook: %+v", err)
	}

	// The following function ensures to atomically create or get the certificate from Kubernetes
	// secret API store. Multiple instances should end up with the same webhookHandlerCert after this function executed.
	crdConverterwebhookHandlerCert, err = providers.GetCertificateFromSecret(osmNamespace, constants.ConversionWebhookCertificateSecretName, crdConverterwebhookHandlerCert, kubeClient)
	if err != nil {
		return errors.Errorf("Error fetching webhook certificate from k8s secret: %s", err)
	}

	wh := converterWebhook{
		config: config,
		cert:   crdConverterwebhookHandlerCert,
	}

	// Start the MutatingWebhook web server
	go wh.run(stop)

	log.Info().Msgf("TEST updating CA bundle of crd")
	// update the crd webhook config with the caBundle
	if err = updateCrdWebhookCABundle(crdConverterwebhookHandlerCert, crdClient, osmNamespace); err != nil {
		return errors.Errorf("Error update crd webook ca bundle %v ", err)
	}

	return nil
}

func (wh *converterWebhook) run(stop <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := http.NewServeMux()

	mux.HandleFunc(WebhookHealthPath, healthHandler)

	// We know that the events arriving at this handler are CREATE POD only
	// because of the specifics of MutatingWebhookConfiguration template in this repository.
	// TODO (snchh): update handler logic to have actual conversion stratergy
	mux.HandleFunc(webhookMeshConfigConvert, serveExampleConvert)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", wh.config.ListenPort),
		Handler: mux,
	}

	log.Info().Msgf("TEST Starting converter webhook server on port: %v", wh.config.ListenPort)
	go func() {
		// Generate a key pair from your pem-encoded cert and key ([]byte).
		cert, err := tls.X509KeyPair(wh.cert.GetCertificateChain(), wh.cert.GetPrivateKey())
		if err != nil {
			log.Error().Err(err).Msg("Error parsing webhook certificate")
			return
		}

		// #nosec G402
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Error().Err(err).Msg("crd conversion webhook HTTP server failed to start")
			return
		}
	}()

	// Wait on exit signals
	<-stop

	// Stop the server
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Error shutting down sidecar-injection webhook HTTP server")
	} else {
		log.Info().Msg("TEST Done shutting down conversion webhook HTTP server")
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Health OK")); err != nil {
		log.Error().Err(err).Msg("Error writing bytes for conversion webhook health check handler")
	}
}

func updateCrdWebhookCABundle(cert certificate.Certificater, crdClient apiclient.ApiextensionsV1Interface, osmNamespace string) error {

	crd, err := crdClient.CustomResourceDefinitions().Get(context.Background(), "meshconfigs.config.openservicemesh.io", metav1.GetOptions{})
	if err != nil {
		return err
	}

	meshConfigCrdPath := "/meshconfigconvert"
	crd.Spec.Conversion = &apiv1.CustomResourceConversion{
		Strategy: apiv1.WebhookConverter,
		Webhook: &apiv1.WebhookConversion{
			ClientConfig: &apiv1.WebhookClientConfig{
				Service: &apiv1.ServiceReference{
					Namespace: osmNamespace,
					Name:      crdConvertedServiceName,
					Path:      &meshConfigCrdPath,
				},
				CABundle: cert.GetCertificateChain(),
			},
			ConversionReviewVersions: []string{"v1alpha2", "v1alpha1", "v1beta1"},
		},
	}
	/*patchJSON, err := json.Marshal(getPartialConversionWebhookConfiguration(cert, osmNamespace))
	if err != nil {
		return err
	}

	if _, err = crd.Patch(context.Background(), "meshconfigs.config.openservicemesh.io", types.StrategicMergePatchType, patchJSON, metav1.PatchOptions{}); err != nil {
		log.Error().Err(err).Msgf("Error updating CA Bundle for conversion webhook %s", "meshconfigs.config.openservicemesh.io")
		return err
	}*/

	if _, err = crdClient.CustomResourceDefinitions().Update(context.Background(), crd, metav1.UpdateOptions{}); err != nil {
		log.Error().Err(err).Msgf("Error updating CA Bundle for conversion webhook %s", "meshconfigs.config.openservicemesh.io")
		return err
	}

	log.Info().Msgf("Finished updating CA Bundle for conversion webhook %s", "meshconfigs.config.openservicemesh.io")
	return nil
}

// getPartialMutatingWebhookConfiguration returns only the portion of the MutatingWebhookConfiguration that needs to be updated.
func getPartialConversionWebhookConfiguration(cert certificate.Certificater, osmNamespace string) apiv1.CustomResourceDefinition {
	meshConfigCrdPath := "/meshconfigconvert"
	return apiv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "meshconfigs.config.openservicemesh.io",
		},
		Spec: apiv1.CustomResourceDefinitionSpec{
			Conversion: &apiv1.CustomResourceConversion{
				Strategy: apiv1.WebhookConverter,
				Webhook: &apiv1.WebhookConversion{
					ClientConfig: &apiv1.WebhookClientConfig{
						Service: &apiv1.ServiceReference{
							Namespace: osmNamespace,
							Name:      crdConvertedServiceName,
							Path:      &meshConfigCrdPath,
						},
						CABundle: cert.GetCertificateChain(),
					},
					ConversionReviewVersions: []string{"v1alpha2", "v1alpha1", "v1beta1"},
				},
			},
		},
	}
}
