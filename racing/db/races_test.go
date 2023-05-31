package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

func getTestData() [][]interface{} {
	var mockData = [][]interface{}{
		// id, meeting_id, name, number, visible, advertised_start_time
		{1, 1, "test1", 1, 0, "2021-09-01T00:00:00Z"},
		{2, 3, "test2", 2, 1, "2021-09-02T00:00:00Z"},
		{3, 2, "test3", 3, 1, time.Now().AddDate(0, 0, -1).Format(time.RFC3339)},
		{4, 1, "test4", 4, 1, time.Now().AddDate(0, 0, 2).Format(time.RFC3339)},
	}
	return mockData
}

func seed(db *sql.DB) error {
	mockData := getTestData()
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS races (id INTEGER PRIMARY KEY, meeting_id INTEGER, name TEXT, number INTEGER, visible INTEGER, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}
	for _, v := range mockData {
		statement, err = db.Prepare(`INSERT OR IGNORE INTO races(id, meeting_id, name, number, visible, advertised_start_time) VALUES (?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(v...)
		}
	}
	return err
}

func TestRacesRepoList(t *testing.T) {
	// Create a test database connection
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Initialize the races repository
	repo := NewRacesRepo(db)
	err = seed(db)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test filter
	filter := &racing.ListRacesRequestFilter{
		Visible: proto.Bool(true),
		OrderBy: "advertised_start_time desc",
	}

	// Call the List method
	races, err := repo.List(filter)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the expected number of races
	expectedCount := 3
	assert.Len(t, races, expectedCount, "Unexpected number of races")
	assert.Equal(t, "test4", races[0].Name, "Unexpected name")
	assert.Equal(t, racing.RaceStatus_OPEN.Enum().String(), races[0].Status.String(), "Unexpected open status")
	assert.Equal(t, racing.RaceStatus_CLOSED.Enum().String(), races[1].Status.String(), "Unexpected close status")
}
