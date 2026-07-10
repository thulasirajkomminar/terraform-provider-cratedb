package provider

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
)

func TestLoggingRoundTripper(t *testing.T) {
	const (
		apiKey    = "test-api-key"
		apiSecret = "super-secret-value"
		respBody  = `{"token":"super-secret-value","name":"example"}`
	)

	var wireAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wireAuthHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer server.Close()

	var logOutput bytes.Buffer
	ctx := tflogtest.RootLogger(context.Background(), &logOutput)

	client := &http.Client{Transport: newLoggingRoundTripper(apiKey, apiSecret, http.DefaultTransport)}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/api/v2/organizations/", nil)
	if err != nil {
		t.Fatalf("building request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("executing request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}

	// The wire request must carry the real, unmasked credentials.
	expectedUser, expectedPass, ok := (&http.Request{Header: http.Header{"Authorization": []string{wireAuthHeader}}}).BasicAuth()
	if !ok || expectedUser != apiKey || expectedPass != apiSecret {
		t.Errorf("wire request is missing the real basic auth credentials, got header %q", wireAuthHeader)
	}

	// The caller must receive the unmasked response body.
	if string(body) != respBody {
		t.Errorf("caller received modified body: got %q, want %q", string(body), respBody)
	}

	logs := logOutput.String()

	// The HTTP transaction must be logged.
	if !strings.Contains(logs, "/api/v2/organizations/") {
		t.Errorf("logs do not contain the request transaction:\n%s", logs)
	}
	if !strings.Contains(logs, "Sending HTTP Request") || !strings.Contains(logs, "Received HTTP Response") {
		t.Errorf("logs do not contain request/response entries:\n%s", logs)
	}

	// The response body must be logged (proving the no-credential assertions
	// below are meaningful), with the secret redacted.
	if !strings.Contains(logs, "example") {
		t.Errorf("logs do not contain the response body:\n%s", logs)
	}

	// The logs must never contain the credentials, even though the response
	// body echoed the secret.
	if strings.Contains(logs, apiSecret) {
		t.Errorf("logs leak the API secret:\n%s", logs)
	}
	if strings.Contains(logs, wireAuthHeader) {
		t.Errorf("logs leak the Authorization header value:\n%s", logs)
	}
}
