package get

import (
	"database/sql"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/quiz"
)

// GetFeacherFromDB ruft die Fächer aus der Datenbank ab und gibt sie als Slice zurück.
func GetFeacherFromDB(db *sql.DB) ([]quizstructs.Schulfach, error) {
	// Führt eine sql-Abfrage aus, um die Fächer zu erhalten.
	rows, err := db.Query(`
		SELECT
			fach
		FROM feacher;
	`)
	if err != nil {
		return nil, err
	}
	// Schließt die Rows am Ende der Funktion.
	defer rows.Close()

	var feacher []quizstructs.Schulfach
	// Iteriert über die erhaltenen Rows.
	for rows.Next() {
		var fach string
		if err := rows.Scan(&fach); err != nil {
			return nil, err
		}

		// Überprüft, ob das Fach bereits in der Liste enthalten ist.
		found := false
		for i := range feacher {
			if feacher[i].Schulfach == fach {
				found = true
				break
			}
		}

		// Wenn das Fach nicht gefunden wurde, füge es zur Liste hinzu.
		if !found {
			feacher = append(feacher, quizstructs.Schulfach{
				Schulfach: fach,
			})
		}
	}

	// Überprüft auf Fehler während des Durchlaufs der Zeilen.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Gibt die Liste der Fächer zurück.
	return feacher, nil
}
