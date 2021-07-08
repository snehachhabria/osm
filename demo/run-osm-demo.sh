#!/bin/bash

set -aueo pipefail

if [ ! -f .env ]; then
    echo -e "\nThere is no .env file in the root of this repository."
    echo -e "Copy the values from .env.example into .env."
    echo -e "Modify the values in .env to match your setup.\n"
    echo -e "    cat .env.example > .env\n\n"
    exit 1
fi

# shellcheck disable=SC1091
source .env

# Set meaningful defaults for env vars we expect from .env
CI="${CI:-false}"  # This is set to true by Github Actions
MESH_NAME="${MESH_NAME:-osm}"
K8S_NAMESPACE="${K8S_NAMESPACE:-osm-system}"
BOOKBUYER_NAMESPACE="${BOOKBUYER_NAMESPACE:-bookbuyer}"
BOOKSTORE_NAMESPACE="${BOOKSTORE_NAMESPACE:-bookstore}"
BOOKTHIEF_NAMESPACE="${BOOKTHIEF_NAMESPACE:-bookthief}"
BOOKWAREHOUSE_NAMESPACE="${BOOKWAREHOUSE_NAMESPACE:-bookwarehouse}"
CERT_MANAGER="${CERT_MANAGER:-tresor}"
CTR_REGISTRY="${CTR_REGISTRY:-localhost:5000}"
CTR_REGISTRY_CREDS_NAME="${CTR_REGISTRY_CREDS_NAME:-acr-creds}"
DEPLOY_TRAFFIC_SPLIT="${DEPLOY_TRAFFIC_SPLIT:-true}"
CTR_TAG="${CTR_TAG:-$(git rev-parse HEAD)}"
IMAGE_PULL_POLICY="${IMAGE_PULL_POLICY:-Always}"
ENABLE_DEBUG_SERVER="${ENABLE_DEBUG_SERVER:-true}"
ENABLE_EGRESS="${ENABLE_EGRESS:-false}"
DEPLOY_GRAFANA="${DEPLOY_GRAFANA:-false}"
DEPLOY_JAEGER="${DEPLOY_JAEGER:-false}"
ENABLE_FLUENTBIT="${ENABLE_FLUENTBIT:-false}"
DEPLOY_PROMETHEUS="${DEPLOY_PROMETHEUS:-false}"
DEPLOY_WITH_SAME_SA="${DEPLOY_WITH_SAME_SA:-false}"
ENVOY_LOG_LEVEL="${ENVOY_LOG_LEVEL:-debug}"
DEPLOY_ON_OPENSHIFT="${DEPLOY_ON_OPENSHIFT:-false}"

# For any additional installation arguments. Used heavily in CI.
optionalInstallArgs=$*

exit_error() {
    error="$1"
    echo "$error"
    exit 1
}

# Check if Docker daemon is running
docker info > /dev/null || { echo "Docker daemon is not running"; exit 1; }

if [ "$CERT_MANAGER" = "vault" ]; then
  # shellcheck disable=SC2086
  bin/osm install \
      --osm-namespace "$K8S_NAMESPACE" \
      --mesh-name "$MESH_NAME" \
      --set=OpenServiceMesh.certificateManager="$CERT_MANAGER" \
      --set=OpenServiceMesh.vault.host="$VAULT_HOST" \
      --set=OpenServiceMesh.vault.token="$VAULT_TOKEN" \
      --set=OpenServiceMesh.vault.protocol="$VAULT_PROTOCOL" \
      --set=OpenServiceMesh.image.registry="$CTR_REGISTRY" \
      --set=OpenServiceMesh.imagePullSecrets[0].name="$CTR_REGISTRY_CREDS_NAME" \
      --set=OpenServiceMesh.image.tag="$CTR_TAG" \
      --set=OpenServiceMesh.image.pullPolicy="$IMAGE_PULL_POLICY" \
      --set=OpenServiceMesh.enableDebugServer="$ENABLE_DEBUG_SERVER" \
      --set=OpenServiceMesh.enableEgress="$ENABLE_EGRESS" \
      --set=OpenServiceMesh.deployGrafana="$DEPLOY_GRAFANA" \
      --set=OpenServiceMesh.deployJaeger="$DEPLOY_JAEGER" \
      --set=OpenServiceMesh.enableFluentbit="$ENABLE_FLUENTBIT" \
      --set=OpenServiceMesh.deployPrometheus="$DEPLOY_PROMETHEUS" \
      --set=OpenServiceMesh.envoyLogLevel="$ENVOY_LOG_LEVEL" \
      --set=OpenServiceMesh.controllerLogLevel="trace" \
      --timeout=90s \
      $optionalInstallArgs
else
  # shellcheck disable=SC2086
  bin/osm install \
      --osm-namespace "$K8S_NAMESPACE" \
      --mesh-name "$MESH_NAME" \
      --set=OpenServiceMesh.certificateManager="$CERT_MANAGER" \
      --set=OpenServiceMesh.image.registry="$CTR_REGISTRY" \
      --set=OpenServiceMesh.imagePullSecrets[0].name="$CTR_REGISTRY_CREDS_NAME" \
      --set=OpenServiceMesh.image.tag="$CTR_TAG" \
      --set=OpenServiceMesh.image.pullPolicy="$IMAGE_PULL_POLICY" \
      --set=OpenServiceMesh.enableDebugServer="$ENABLE_DEBUG_SERVER" \
      --set=OpenServiceMesh.enableEgress="$ENABLE_EGRESS" \
      --set=OpenServiceMesh.deployGrafana="$DEPLOY_GRAFANA" \
      --set=OpenServiceMesh.deployJaeger="$DEPLOY_JAEGER" \
      --set=OpenServiceMesh.enableFluentbit="$ENABLE_FLUENTBIT" \
      --set=OpenServiceMesh.deployPrometheus="$DEPLOY_PROMETHEUS" \
      --set=OpenServiceMesh.envoyLogLevel="$ENVOY_LOG_LEVEL" \
      --set=OpenServiceMesh.controllerLogLevel="trace" \
      --timeout=90s \
      $optionalInstallArgs
fi

./demo/configure-app-namespaces.sh

./demo/deploy-apps.sh

# Apply SMI policies
if [ "$DEPLOY_TRAFFIC_SPLIT" = "true" ]; then
    ./demo/deploy-traffic-split.sh
fi

./demo/deploy-traffic-specs.sh

if [ "$DEPLOY_WITH_SAME_SA" = "true" ]; then
    ./demo/deploy-traffic-target-with-same-sa.sh
else
    ./demo/deploy-traffic-target.sh
fi

if [[ "$CI" != "true" ]]; then
    watch -n5 "printf \"Namespace ${K8S_NAMESPACE}:\n\"; kubectl get pods -n ${K8S_NAMESPACE} -o wide; printf \"\n\n\"; printf \"Namespace ${BOOKBUYER_NAMESPACE}:\n\"; kubectl get pods -n ${BOOKBUYER_NAMESPACE} -o wide; printf \"\n\n\"; printf \"Namespace ${BOOKSTORE_NAMESPACE}:\n\"; kubectl get pods -n ${BOOKSTORE_NAMESPACE} -o wide; printf \"\n\n\"; printf \"Namespace ${BOOKTHIEF_NAMESPACE}:\n\"; kubectl get pods -n ${BOOKTHIEF_NAMESPACE} -o wide; printf \"\n\n\"; printf \"Namespace ${BOOKWAREHOUSE_NAMESPACE}:\n\"; kubectl get pods -n ${BOOKWAREHOUSE_NAMESPACE} -o wide"
fi
