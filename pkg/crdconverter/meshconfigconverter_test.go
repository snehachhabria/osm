package crdconverter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tassert "github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func TestConverter(t *testing.T) {
	assert := tassert.New(t)

	sampleObj := `kind: ConversionReview
apiVersion: config.openservicemesh.io/v1alpha1
request:
  uid: 0000-0000-0000-0000
  desiredAPIVersion: config.openservicemesh.io/v1alpha2
  objects:
  - apiVersion: config.openservicemesh.io/v1alpha1
    kind: MeshConfig
    metadata:
      name: preset-mesh-config-object
    spec:
      cronSpec: "* * * * */5"
      host: "abc.com"
      image: my-awesome-cron-image
      featureFlags:
       enableEgressPolicy: true
       enableMulticlusterMode: false
       enableWASMStats: true
`
	// First try json, it should fail as the data is taml
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/convert", strings.NewReader(sampleObj))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/json")
	serveExampleConvert(response, request)
	convertReview := v1beta1.ConversionReview{}
	scheme := runtime.NewScheme()
	jsonSerializer := json.NewSerializer(json.DefaultMetaFactory, scheme, scheme, false)
	if _, _, err := jsonSerializer.Decode(response.Body.Bytes(), nil, &convertReview); err != nil {
		t.Fatal(err)
	}
	if convertReview.Response.Result.Status != v1.StatusFailure {
		t.Fatalf("expected the operation to fail when yaml is provided with json header")
	} else if !strings.Contains(convertReview.Response.Result.Message, "json parse error") {
		t.Fatalf("expected to fail on json parser, but it failed with: %v", convertReview.Response.Result.Message)
	}

	// Now try yaml, and it should successfully convert
	response = httptest.NewRecorder()
	request, err = http.NewRequest("POST", "/convert", strings.NewReader(sampleObj))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/yaml")
	serveExampleConvert(response, request)
	convertReview = v1beta1.ConversionReview{}
	yamlSerializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme, scheme)
	if _, _, err := yamlSerializer.Decode(response.Body.Bytes(), nil, &convertReview); err != nil {
		t.Fatalf("cannot decode data: \n %v\n Error: %v", response.Body, err)
	}
	if convertReview.Response.Result.Status != v1.StatusSuccess {
		t.Fatalf("cr conversion failed: %v", convertReview.Response)
	}
	convertedObj := unstructured.Unstructured{}
	if _, _, err := yamlSerializer.Decode(convertReview.Response.ConvertedObjects[0].Raw, nil, &convertedObj); err != nil {
		t.Fatal(err)
	}
	if e, a := "config.openservicemesh.io/v1alpha2", convertedObj.GetAPIVersion(); e != a {
		t.Errorf("expected= %v, actual= %v", e, a)
	}

	actualSpec := convertedObj.Object["spec"]
	assert.NotEmpty(actualSpec)
}
