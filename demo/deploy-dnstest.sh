set -aueo pipefail

# shellcheck disable=SC1091
source .env

echo -e "Deploy DnsTest Service Account"
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: dnsutils
  namespace: bookbuyer
EOF

echo -e "Deploy DnsTest Deployment"
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dnsutils
  namespace: bookbuyer
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dnsutils
      version: v1
  template:
    metadata:
      labels:
        app: dnsutils
        version: v1
    spec:
      serviceAccountName: dnsutils
      nodeSelector:
        kubernetes.io/arch: amd64
        kubernetes.io/os: linux
      containers:
        # Main container with APP
        - name: dnsutils
          image: gcr.io/kubernetes-e2e-test-images/dnsutils:1.3
          imagePullPolicy: Always
          command:
            - sleep
            - "3600"

      imagePullSecrets:
        - name: "$CTR_REGISTRY_CREDS_NAME"
EOF