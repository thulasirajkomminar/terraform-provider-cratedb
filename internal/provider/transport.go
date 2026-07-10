package provider

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// authTransport injects the credentials into every outgoing request. It MUST
// sit below the logging transport in the chain so the Authorization header is
// added after the request has been logged and therefore never appears in logs.
type authTransport struct {
	apiKey    string
	apiSecret string
	base      http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.SetBasicAuth(t.apiKey, t.apiSecret)
	return t.base.RoundTrip(req)
}

// maskingTransport sits above the logging transport and injects log masking
// rules into the request context so any credential value that shows up in a
// logged request or response (e.g. echoed in a body) is redacted by tflog.
type maskingTransport struct {
	masks []*regexp.Regexp
	base  http.RoundTripper
}

func (t *maskingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if len(t.masks) > 0 {
		ctx = tflog.MaskAllFieldValuesRegexes(ctx, t.masks...)
		ctx = tflog.MaskMessageRegexes(ctx, t.masks...)
	}
	return t.base.RoundTrip(req.WithContext(ctx))
}

// credentialMasks builds the redaction patterns for every representation of
// the credentials that could appear on the wire, including the base64-encoded
// basic auth token.
func credentialMasks(apiKey, apiSecret string) []*regexp.Regexp {
	var masks []*regexp.Regexp
	for _, secret := range []string{
		apiSecret,
		base64.StdEncoding.EncodeToString([]byte(apiKey + ":" + apiSecret)),
	} {
		if secret != "" {
			masks = append(masks, regexp.MustCompile(regexp.QuoteMeta(secret)))
		}
	}
	return masks
}

// newAPIHTTPClient builds the retrying HTTP client used by the CrateDB API
// client. Each retry attempt flows through the transport chain
// masking -> logging -> auth -> base, so with TF_LOG=DEBUG every request and
// response (per attempt) is logged without ever exposing credentials.
func newAPIHTTPClient(apiKey, apiSecret string) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.RetryMax = 3
	retryClient.HTTPClient.Transport = newLoggingRoundTripper(apiKey, apiSecret, http.DefaultTransport)
	return retryClient.StandardClient()
}

// newLoggingRoundTripper wires the masking -> logging -> auth -> base chain
// around the given base transport.
func newLoggingRoundTripper(apiKey, apiSecret string, base http.RoundTripper) http.RoundTripper {
	return &maskingTransport{
		masks: credentialMasks(apiKey, apiSecret),
		base: logging.NewLoggingHTTPTransport(&authTransport{
			apiKey:    apiKey,
			apiSecret: apiSecret,
			base:      base,
		}),
	}
}
