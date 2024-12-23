package diddht

import (
	"context"
	"net/http"

	"github.com/decentralized-identity/web5-go/dids/did"
	"github.com/decentralized-identity/web5-go/dids/didcore"
	"github.com/decentralized-identity/web5-go/dids/diddht/internal/dns"
	"github.com/decentralized-identity/web5-go/dids/diddht/internal/pkarr"
	"github.com/tv42/zbase32"
)

// DefaultResolver uses the default Pkarr gateway client: https://diddht.tbddev.org
func DefaultResolver() *Resolver {
	return &Resolver{
		relay: getDefaultGateway(),
	}
}

// Resolver is a client for resolving DIDs using the DHT network.
type Resolver struct {
	relay gateway
}

// NewResolver creates a new Resolver instance with the given relay and HTTP client.
// TODO make this relay an option and use default relay if not provided
func NewResolver(relayURL string, client *http.Client) *Resolver {
	pkarrRelay := pkarr.NewClient(relayURL, client)
	return &Resolver{
		relay: pkarrRelay,
	}
}

// Resolve resolves a DID using the DHT method
func (r *Resolver) Resolve(uri string) (didcore.ResolutionResult, error) {
	return r.ResolveWithContext(context.Background(), uri)
}

// ResolveWithContext resolves a DID using the DHT method. This is the context aware version of Resolve.
func (r *Resolver) ResolveWithContext(ctx context.Context, uri string) (didcore.ResolutionResult, error) {

	// 1. Parse URI and make sure it's a DHT method
	did, err := did.Parse(uri)
	if err != nil {
		// TODO log err
		return didcore.ResolutionResultWithError("invalidDid"), didcore.ResolutionError{Code: "invalidDid"}
	}

	if did.Method != "dht" {
		return didcore.ResolutionResultWithError("methodNotSupported"), didcore.ResolutionError{Code: "methodNotSupported"}
	}

	// 2. ensure did ID is zbase32
	identifier, err := zbase32.DecodeString(did.ID)
	if err != nil {
		// TODO log err
		return didcore.ResolutionResultWithError("invalidPublicKey"), didcore.ResolutionError{Code: "invalidPublicKey"}
	}

	if len(identifier) == 0 {
		// return nil, fmt.Errorf("no bytes decoded from zbase32 identifier %s", did.ID)
		// TODO log err
		return didcore.ResolutionResultWithError("invalidPublicKey"), didcore.ResolutionError{Code: "invalidPublicKey"}
	}

	// 3. fetch from the relay
	bep44Message, err := r.relay.FetchWithContext(ctx, did.ID)
	if err != nil {
		// TODO log err
		return didcore.ResolutionResultWithError("notFound"), didcore.ResolutionError{Code: "notFound"}
	}

	// get the dns payload from the bep44 message
	bep44MessagePayload := bep44Message.V
	document, err := dns.UnmarshalDIDDocument(bep44MessagePayload)
	if err != nil {
		// TODO log err
		return didcore.ResolutionResultWithError("invalidDid"), didcore.ResolutionError{Code: "invalidDid"}
	}

	return didcore.ResolutionResultWithDocument(*document), nil
}
