<p align="center">
  <img src="https://raw.githubusercontent.com/cert-manager/cert-manager/d53c0b9270f8cd90d908460d69502694e1838f5f/logo/logo-small.png" height="256" width="256" alt="cert-manager project logo" />
</p>

# cert-manager - ExternalDNS webhook

This repo allows [cert-manager](https://github.com/cert-manager/cert-manager) to use [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) to handle ACME challenges. 

> [!NOTE]
> This repository started as a fork of
> [lion7/cert-manager-webhook-external-dns](https://github.com/lion7/cert-manager-webhook-external-dns)
> by [Gerard de Leeuw](https://github.com/lion7) and is now maintained by [youniqx](https://github.com/youniqx)
> as an independent, standalone project. Many thanks to the original author for the groundwork.
> The two projects are **not** interchangeable: this fork uses its own API group
> (`cert-manager-webhook-external-dns.youniqx.com`), its own container image and Helm chart location.

## Requirements 

### ExternalDNS

The default configuration of ExternalDNS needs altering for this integration to function:
- TXT records are not managed by ExternalDNS by default, it requires an extra flag 
- The DNSEndpoint CRD is not enabled by default 

If you are deploying with the [official Helm chart](https://artifacthub.io/packages/helm/external-dns/external-dns) you can accomplish this by including this in your values file:

```yaml
managedRecordTypes: 
  - A      # ┐
  - AAAA   # ├ Default values
  - CNAME  # ┘
  - TXT    # ─ New value

sources:
  - service # ┬ Default values
  - ingress # ┘
  - crd     # ─ New value
```

### cert-manager

Any supported version of cert-manager supports DNS webhooks, for documentation on installing cert-manager see the [official documentation](https://cert-manager.io/docs/installation/)

## Installing

### Using Helm

The webhook can be installed using Helm:

```bash
helm install external-dns-webhook \
  oci://ghcr.io/youniqx/cert-manager-webhook-external-dns/helm/external-dns-webhook:<tag> \
  --namespace cert-manager
```

### Configuration

The chart is built on the [bjw-s common library](https://bjw-s-labs.github.io/helm-charts/docs/common-library/), so all of its
values (controllers, service, rbac, persistence, ...) are available for customization. The default values should work for most
installations. You can customize the installation by creating a values file:

```yaml
# values.yaml
controllers:
  webhook:
    replicas: 1
    containers:
      app:
        image:
          repository: ghcr.io/youniqx/cert-manager-webhook-external-dns
          # Defaults to the chart's appVersion
          tag: "{{ .Chart.AppVersion }}"
          pullPolicy: IfNotPresent
        resources:
          limits:
            cpu: 100m
            memory: 128Mi
          requests:
            cpu: 100m
            memory: 128Mi

# cert-manager configuration
certManager:
  namespace: cert-manager
  serviceAccountName: cert-manager
```

Then install with:

```bash
helm install external-dns-webhook \
  oci://ghcr.io/youniqx/cert-manager-webhook-external-dns/helm/external-dns-webhook:<tag> \
  --namespace cert-manager \
  -f values.yaml
```

## Supply chain security

Container images are multi-platform (`linux/amd64`, `linux/arm64`), built from a [distroless](https://github.com/GoogleContainerTools/distroless)
static base image with all base images pinned by digest. Every release image and Helm chart is signed with
[cosign](https://docs.sigstore.dev/) keyless signing and ships with SLSA build provenance attestations and a CycloneDX SBOM.

Verify the image signature:

```bash
# Image
cosign verify ghcr.io/youniqx/cert-manager-webhook-external-dns:<tag> \
  --certificate-identity-regexp 'https://github.com/youniqx/cert-manager-webhook-external-dns/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com

# Helm chart
cosign verify ghcr.io/youniqx/cert-manager-webhook-external-dns/helm/external-dns-webhook:<tag> \
  --certificate-identity-regexp 'https://github.com/youniqx/cert-manager-webhook-external-dns/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```

Verify the build provenance (image or chart) via the GitHub CLI:

```bash
# Image
gh attestation verify oci://ghcr.io/youniqx/cert-manager-webhook-external-dns:<tag> \
  --repo youniqx/cert-manager-webhook-external-dns

# Helm chart
gh attestation verify oci://ghcr.io/youniqx/cert-manager-webhook-external-dns/helm/external-dns-webhook:<tag> \
  --repo youniqx/cert-manager-webhook-external-dns
```

Inspect the CycloneDX SBOM attached to the image (also uploaded as a release asset):

```bash
cosign verify-attestation ghcr.io/youniqx/cert-manager-webhook-external-dns:<tag> \
  --type cyclonedx \
  --certificate-identity-regexp 'https://github.com/youniqx/cert-manager-webhook-external-dns/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  | jq -r '.payload' | base64 -d | jq '.predicate'
```

## Usage

To configure an issuer to use ExternalDNS you just specify the group and solver name within the Issuer or ClusterIssuer config:

```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: example-issuer
spec:
  acme:
   ...
    solvers:
    - dns01:
        webhook:
          groupName: cert-manager-webhook-external-dns.youniqx.com
          solverName: external-dns
```

## Running the test suite

All DNS providers **must** run the DNS01 provider conformance testing suite,
else they will have undetermined behaviour when used with cert-manager.

**It is essential that you configure and run the test suite when creating a
DNS01 webhook.**

A test file has been provided in [solver_test.go](pkg/externaldnssolver/solver_test.go), which runs the official cert-manager
DNS01 conformance suite against this webhook. The required Kubernetes test binaries are downloaded automatically
(see [TESTING.md](TESTING.md) for details):

```bash
make test
```

The test suite includes a mock DNS server with external-dns controller simulation to provide comprehensive testing without requiring external DNS provider credentials.
