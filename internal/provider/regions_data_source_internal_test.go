package provider

import (
	"encoding/json"
	"testing"
	"time"
)

func TestApiRegionTolerantTimestamps(t *testing.T) {
	// Exact shape returned by the live regions endpoint: last_seen has no
	// timezone offset, dc timestamps do.
	payload := `[
		{
			"name": "aks1.westeurope.azure",
			"description": "West Europe (Azure)",
			"deprecated": false,
			"is_edge_region": false,
			"status": "UP",
			"upgrade_available": false,
			"last_seen": "2026-07-10T10:41:02.983000",
			"dc": {"created": "2024-10-08T10:00:00+00:00", "modified": "2024-10-08T11:00:00Z"}
		}
	]`

	var regions []apiRegion
	if err := json.Unmarshal([]byte(payload), &regions); err != nil {
		t.Fatalf("unmarshalling regions payload: %v", err)
	}

	region := regions[0]
	if region.LastSeen == nil || !region.LastSeen.Equal(time.Date(2026, 7, 10, 10, 41, 2, 983000000, time.UTC)) {
		t.Errorf("last_seen not parsed as UTC: %v", region.LastSeen)
	}
	if region.Dc == nil || region.Dc.Created == nil || !region.Dc.Created.Equal(time.Date(2024, 10, 8, 10, 0, 0, 0, time.UTC)) {
		t.Errorf("dc.created not parsed: %v", region.Dc)
	}

	var invalid apiTime
	if err := json.Unmarshal([]byte(`"not-a-timestamp"`), &invalid); err == nil {
		t.Error("expected an error for an invalid timestamp")
	}
}
