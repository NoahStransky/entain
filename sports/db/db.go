package db

import (
	"fmt"
	"time"

	"syreclabs.com/go/faker"
)

func (s *sportsRepo) seed() error {
	statement, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS sports (id INTEGER PRIMARY KEY, meeting_id INTEGER, name TEXT, visible INTEGER, advertised_start_time DATETIME, description TEXT)`)
	if err == nil {
		fmt.Println("create table")
		_, err = statement.Exec()
	}
	fmt.Println("error11 ", err)

	for i := 1; i <= 100; i++ {
		statement, err = s.db.Prepare(`INSERT OR IGNORE INTO sports(id, meeting_id, name, visible, advertised_start_time, description) VALUES (?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Number().Between(1, 10),
				faker.Team().Name(),
				faker.Number().Between(0, 1),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
				faker.Lorem().Sentence(10),
			)
		}
	}
	fmt.Println("error ", err)

	return err
}
