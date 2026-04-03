package mlbapi

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

const integrationEnvVar = "MLBAPI_RUN_INTEGRATION"

func requireIntegration(t *testing.T) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping live MLB API test in short mode")
	}
	if os.Getenv(integrationEnvVar) != "1" {
		t.Skipf("set %s=1 to run live MLB API tests", integrationEnvVar)
	}
}

func newLiveClient() (*Client, *statusRecorderTransport) {
	recorder := &statusRecorderTransport{}
	httpClient := &http.Client{
		Timeout:   20 * time.Second,
		Transport: recorder,
	}

	return NewClient(WithHTTPClient(httpClient)), recorder
}

func assertStatuses200(t *testing.T, statuses []int) {
	t.Helper()

	if len(statuses) == 0 {
		t.Fatal("expected at least one live HTTP response status")
	}
	for _, status := range statuses {
		if status != http.StatusOK {
			t.Fatalf("expected all live responses to be 200, got %v", statuses)
		}
	}
}

func TestLiveMetaPositionsReturns200(t *testing.T) {
	requireIntegration(t)

	client, recorder := newLiveClient()
	payload, err := client.Meta(context.Background(), "positions")
	if err != nil {
		t.Fatalf("Meta returned error: %v", err)
	}

	assertStatuses200(t, recorder.Statuses())
	if len(nestedSlice(payload, "positions")) == 0 {
		t.Fatalf("expected positions payload to be non-empty: %+v", payload)
	}
}

func TestLiveTeamsEndpointReturns200(t *testing.T) {
	requireIntegration(t)

	client, recorder := newLiveClient()
	payload, err := client.Get(context.Background(), "teams", Params{"sportIds": 1})
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	assertStatuses200(t, recorder.Statuses())
	if len(nestedSlice(payload, "teams")) < 30 {
		t.Fatalf("expected MLB teams payload, got %+v", payload)
	}
}

func TestLiveLookupTeamReturns200(t *testing.T) {
	requireIntegration(t)

	client, recorder := newLiveClient()
	teams, err := client.LookupTeam(context.Background(), "dodgers", LookupTeamOptions{
		SportIDs: "1",
		Season:   currentYear(),
	})
	if err != nil {
		t.Fatalf("LookupTeam returned error: %v", err)
	}

	assertStatuses200(t, recorder.Statuses())
	if len(teams) == 0 {
		t.Fatal("expected at least one team match from live lookup")
	}

	foundDodgers := false
	for _, team := range teams {
		if team.ID == TeamLosAngelesDodgers {
			foundDodgers = true
			break
		}
	}
	if !foundDodgers {
		t.Fatalf("expected Dodgers in live lookup results: %+v", teams)
	}
}
