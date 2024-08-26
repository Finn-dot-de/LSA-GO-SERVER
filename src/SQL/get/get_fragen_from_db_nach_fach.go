package get

import (
	"database/sql"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/quiz"
)

// GetFragenFromDBNachFach ruft alle Fragen zu einem bestimmten Thema aus der Datenbank ab und gibt sie zurück.
func GetFragenFromDBNachFach(db *sql.DB, FachName string) ([]quizstructs.Frage, error) {
	// Führt eine sql-Abfrage aus, um die Fragen zu einem bestimmten Fach zu erhalten.
	rows, err := db.Query(`
	SELECT
	    fragen.id,
	    fragen.frage_text,
	    feacher.id,
	    feacher.kurzform,
	    feacher.langform,
	    a.id,
	    a.antwort_text,
	    a.ist_korrekt
	FROM fragen
	    JOIN moegliche_antworten AS a ON a.frage_id = fragen.id
	    JOIN feacher ON fragen.fach_id = feacher.id
	WHERE feacher.kurzform = $1;`, FachName)
	if err != nil {
		return nil, err
	}
	// Schließt die Rows am Ende der Funktion.
	defer rows.Close()

	var fragen []quizstructs.Frage
	// Iteriert über die erhaltenen Rows.
	for rows.Next() {
		var frageID, themaID int
		var frageText, themaName, langform string
		var antwortID sql.NullInt64
		var antwortText sql.NullString
		var istKorrekt sql.NullBool

		// Daten aus der Abfrage in Variablen scannen.
		err := rows.Scan(&frageID, &frageText, &themaID, &themaName, &langform, &antwortID, &antwortText, &istKorrekt)
		if err != nil {
			return nil, err
		}

		// Überprüfen, ob die Frage bereits in der Liste vorhanden ist.
		found := false
		for i := range fragen {
			if fragen[i].ID == frageID {
				found = true
				break
			}
		}

		// Wenn die Frage nicht gefunden wurde, füge sie der Liste hinzu.
		if !found {
			fragen = append(fragen, quizstructs.Frage{
				ID:        frageID,
				FrageText: frageText,
				FachID:    themaID,
				FachName:  themaName,
				Langform:  langform,
				Antworten: []quizstructs.Antwort{},
			})
		}

		// Füge die Antwort der entsprechenden Frage hinzu, falls vorhanden.
		if antwortID.Valid {
			fragen[len(fragen)-1].Antworten = append(fragen[len(fragen)-1].Antworten, quizstructs.Antwort{
				AntwortID:  int(antwortID.Int64),
				Antwort:    antwortText.String,
				IstKorrekt: istKorrekt.Bool,
			})
		}
	}
	// Überprüfen auf Fehler während des Durchlaufs der Zeilen.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Gibt die Liste der Fragen zurück.
	return fragen, nil
}
