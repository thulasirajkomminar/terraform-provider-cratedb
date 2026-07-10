package provider

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
)

// clientFromProviderData converts the provider data passed to a resource or
// data source Configure method into the CrateDB API client. It returns nil
// without appending diagnostics when the provider has not been configured yet.
func clientFromProviderData(providerData any, kind string, diags *diag.Diagnostics) *cratedb.ClientWithResponses {
	if providerData == nil {
		return nil
	}

	client, ok := providerData.(*cratedb.ClientWithResponses)
	if !ok {
		diags.AddError(
			fmt.Sprintf("Unexpected %s Configure Type", kind),
			fmt.Sprintf("Expected *cratedb.ClientWithResponses, got: %T. Please report this issue to the provider developers.", providerData),
		)
		return nil
	}
	return client
}

// maxErrorBodyBytes limits how much of a raw API response body is included in
// diagnostics so a huge (e.g. HTML) response does not flood the output.
const maxErrorBodyBytes = 4096

// apiErrorDetail formats an unexpected API response for a diagnostic. It never
// assumes the body matches a documented error shape: it always reports the
// HTTP status and falls back to the raw response body so proxy or load
// balancer interceptions surface the real problem.
func apiErrorDetail(httpResponse *http.Response, body []byte) string {
	status := "unknown"
	statusCode := 0
	if httpResponse != nil {
		status = httpResponse.Status
		statusCode = httpResponse.StatusCode
	}

	detail := fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", statusCode, status)

	rawBody := strings.TrimSpace(string(body))
	if rawBody == "" {
		return detail + "\nResponse Body: (empty)"
	}
	if len(rawBody) > maxErrorBodyBytes {
		rawBody = rawBody[:maxErrorBodyBytes] + "… (truncated)"
	}
	return detail + "\nResponse Body: " + rawBody
}

// isNotFound reports whether an API response is a confirmed "not found", i.e.
// the remote object no longer exists and the resource should be removed from
// state so Terraform plans a re-create.
func isNotFound(httpResponse *http.Response) bool {
	return httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound
}
