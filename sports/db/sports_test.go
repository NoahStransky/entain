package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

func getTestData() [][]interface{} {
	var mockData = [][]interface{}{
		// id, meeting_id, name, visible, advertised_start_time, description
		{1, 1, "event1", 0, "2021-09-01T00:00:00Z", "test event 1"},
		{2, 3, "event2", 1, "2021-09-02T00:00:00Z", "test event 2"},
		{3, 2, "event3", 1, time.Now().AddDate(0, 0, -1).Format(time.RFC3339), "test event 3"},
		{4, 1, "event4", 1, time.Now().AddDate(0, 0, 2).Format(time.RFC3339), "test event 4"},
	}
	return mockData
}

func seed(db *sql.DB) error {
	mockData := getTestData()
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS sports (id INTEGER PRIMARY KEY, meeting_id INTEGER, name TEXT, visible INTEGER, advertised_start_time DATETIME, description TEXT)`)
	if err == nil {
		_, err = statement.Exec()
	}
	for _, v := range mockData {
		statement, err = db.Prepare(`INSERT OR IGNORE INTO sports(id, meeting_id, name, visible, advertised_start_time, description) VALUES (?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(v...)
		}
	}
	return err
}

func TestSportsRepoList(t *testing.T) {
	// Create a test database connection
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Initialize the sports repository
	repo := NewSportsRepo(db)
	err = seed(db)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test filter
	filter := &sports.ListEventsRequestFilter{
		Visible: proto.Bool(true),
		OrderBy: "advertised_start_time desc",
	}

	// Call the List method
	events, err := repo.List(filter)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the expected number of events
	expectedCount := 3
	assert.Len(t, events, expectedCount, "Unexpected number of events")
	assert.Equal(t, "event4", events[0].Name, "Unexpected name")
}
