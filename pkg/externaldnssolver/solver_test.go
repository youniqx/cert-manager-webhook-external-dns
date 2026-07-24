package externaldnssolver

import (
	"net"
	"testing"
	"time"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
	"github.com/go-logr/logr/testr"
	inttesting "github.com/youniqx/cert-manager-webhook-external-dns/testing"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const DNSPort = "59351"

func TestRunsSuite(t *testing.T) {
	log.SetLogger(testr.New(t))

	solver := inttesting.NewTestSolver(t,
		&ExternalDNSProviderSolver{RequestTimeout: 5 * time.Second},
		net.JoinHostPort("", DNSPort),
		"example.com.",
	)

	fixture := acmetest.NewFixture(solver,
		acmetest.SetResolvedZone("example.com."),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetConfig(map[string]any{}),
		acmetest.SetDNSServer(net.JoinHostPort("127.0.0.1", DNSPort)),
		acmetest.SetUseAuthoritative(false),
		acmetest.SetPropagationLimit(time.Second*30),
	)

	fixture.RunConformance(t)
}
