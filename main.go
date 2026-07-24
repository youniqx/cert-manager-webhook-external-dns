package main

import (
	"time"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/youniqx/cert-manager-webhook-external-dns/pkg/externaldnssolver"
)

const (
	// GroupName is the Kubernetes group name that will be forwarded to this
	// extension-apiserver.
	GroupName = "cert-manager-webhook-external-dns.youniqx.com"

	// RequestTimeout is the timeout for each request
	RequestTimeout = 5 * time.Second
)

func main() {
	cmd.RunWebhookServer(GroupName, &externaldnssolver.ExternalDNSProviderSolver{
		RequestTimeout: RequestTimeout,
	})
}
