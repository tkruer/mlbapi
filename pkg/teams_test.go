package mlbapi

import "testing"

func TestKnownTeamsMetadata(t *testing.T) {
	t.Parallel()

	if len(AllTeams()) != 30 {
		t.Fatalf("expected 30 current MLB teams, got %d", len(AllTeams()))
	}

	dodgers, ok := TeamByID(TeamLosAngelesDodgers)
	if !ok {
		t.Fatal("expected Dodgers metadata to exist")
	}
	if dodgers.Abbreviation != "LAD" || dodgers.TeamCode != "lan" {
		t.Fatalf("unexpected Dodgers metadata: %+v", dodgers)
	}

	athletics, ok := TeamByAbbreviation("ath")
	if !ok {
		t.Fatal("expected Athletics metadata to exist")
	}
	if athletics.ID != TeamAthletics || athletics.LocationName != "Sacramento" {
		t.Fatalf("unexpected Athletics metadata: %+v", athletics)
	}

	byName, ok := TeamByName("New York Yankees")
	if !ok || byName.ID != TeamNewYorkYankees {
		t.Fatalf("expected exact name lookup to resolve Yankees, got %+v ok=%v", byName, ok)
	}
}

func TestTeamIDHelpersAndSorting(t *testing.T) {
	t.Parallel()

	if TeamLosAngelesDodgers.Int() != 119 {
		t.Fatalf("unexpected TeamID integer value: %d", TeamLosAngelesDodgers.Int())
	}
	if TeamLosAngelesDodgers.String() != "Los Angeles Dodgers" {
		t.Fatalf("unexpected TeamID string value: %s", TeamLosAngelesDodgers.String())
	}

	teams := []Team{
		{ID: TeamNewYorkYankees, Name: "New York Yankees"},
		{ID: TeamBostonRedSox, Name: "Boston Red Sox"},
		{ID: TeamLosAngelesDodgers, Name: "Los Angeles Dodgers"},
	}
	SortTeamsByID(teams)

	if teams[0].ID != TeamBostonRedSox || teams[1].ID != TeamLosAngelesDodgers || teams[2].ID != TeamNewYorkYankees {
		t.Fatalf("teams not sorted by ID: %+v", teams)
	}
}
