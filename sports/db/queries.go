package db

const (
	sportsList = "list"
)

func getSportsQueries() map[string]string {
	return map[string]string{
		sportsList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				visible, 
				advertised_start_time,
				description
			FROM sports
		`,
	}
}
