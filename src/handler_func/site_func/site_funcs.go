package site_func

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
	"log"
)

// GetFeacherFromDB ruft die Fächer aus der Datenbank ab und gibt sie als Slice zurück.
func GetFeacherFromDB(db *sql.DB) ([]structs.Schulfach, error) {
	// SQL-Abfrage, um die Fächer abzurufen.
	query := `
		SELECT kurzform
		FROM feacher;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abrufen der Fächer: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var feacher []structs.Schulfach

	// Verarbeitet die Rows und fügt sie zur gegebenen Fächer-Liste hinzu
	for rows.Next() {
		var fach structs.Schulfach
		// Scannt die Row und fügt das Fach hinzu
		err := rows.Scan(&fach.Schulfach)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der Fächer-Daten: %v", err)
		}
		feacher = append(feacher, fach)
	}

	// Überprüft auf fehler während des Durchlaufs der Zeilen
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fehler bei der Verarbeitung der Fächer: %v", err)
	}

	// Gibt die Liste der Fächer zurück
	return feacher, nil
}

// GetLernseiteByID Funktion, um eine Seite basierend auf der ID aus der Datenbank abzurufen
func GetLernseiteByID(db *sql.DB, titel string) (*structs.Lernseite, error) {
	query := `
		SELECT id, titel, daten, benutzer_id, erstellungsdatum
		FROM lernseiten
		WHERE titel = $1;
	`

	row := db.QueryRow(query, titel)

	var lernseite structs.Lernseite
	err := row.Scan(&lernseite.ID, &lernseite.Titel, &lernseite.Text, &lernseite.BenutzerID, &lernseite.Erstellungsdatum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("keine Datei mit dem Titel %s gefunden", titel)
		}
		return nil, fmt.Errorf("fehler beim Abrufen der Datei: %v", err)
	}

	// Konvertiere die Bytes zurück in einen lesbaren String
	lernseite.Text = string(lernseite.Text)

	return &lernseite, err
}

// GetFragenFromDBByFach ruft alle Fragen zu einem bestimmten Thema aus der Datenbank ab und gibt sie zurück.
func GetFragenFromDBByFach(db *sql.DB, FachName string) ([]structs.Frage, error) {
	query := `
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
	WHERE feacher.kurzform = $1;
	`

	rows, err := db.Query(query, FachName)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abrufen der Fragen: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var fragen []structs.Frage
	err = scanFragen(rows, &fragen)
	if err != nil {
		return nil, err
	}

	return fragen, nil
}

// scanFragen verarbeitet die Rows und fügt sie zur gegebenen Fragen-Liste hinzu
func scanFragen(rows *sql.Rows, fragen *[]structs.Frage) error {
	for rows.Next() {
		var frageID, themaID int
		var frageText, themaName, langform string
		var antwortID sql.NullInt64
		var antwortText sql.NullString
		var istKorrekt sql.NullBool

		err := rows.Scan(&frageID, &frageText, &themaID, &themaName, &langform, &antwortID, &antwortText, &istKorrekt)
		if err != nil {
			return fmt.Errorf("fehler beim Scannen der Fragen-Daten: %v", err)
		}

		// Füge Frage und Antwort zur Liste hinzu, wenn sie noch nicht vorhanden ist
		found := false
		for i := range *fragen {
			if (*fragen)[i].ID == frageID {
				found = true
				break
			}
		}

		if !found {
			*fragen = append(*fragen, structs.Frage{
				ID:        frageID,
				FrageText: frageText,
				FachID:    themaID,
				FachName:  themaName,
				Langform:  langform,
				Antworten: []structs.Antwort{},
			})
		}

		if antwortID.Valid {
			(*fragen)[len(*fragen)-1].Antworten = append((*fragen)[len(*fragen)-1].Antworten, structs.Antwort{
				AntwortID:  int(antwortID.Int64),
				Antwort:    antwortText.String,
				IstKorrekt: istKorrekt.Bool,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("fehler bei der Verarbeitung der Fragen: %v", err)
	}

	return nil
}

// CheckIfLernseiteExists prüft, ob eine Seite in der Datenbank existiert
func CheckIfLernseiteExists(db *sql.DB, titel string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM lernseiten WHERE titel=$1)"
	err := db.QueryRow(query, titel).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Überprüfen der Seite: %v", err)
	}
	return exists, nil
}

// SaveLearningSitesByID speichert oder aktualisiert eine Lernseite basierend auf dem Titel.
// Die Funktion konvertiert die Eingabedaten in Byte-Array (bytea), um mit der Datenbank kompatibel zu sein.
func SaveLearningSitesByID(db *sql.DB, titel string, daten interface{}, benutzerID int) error {
	// Daten in []byte umwandeln, basierend auf dem Typ des Eingabewerts.
	var byteData []byte
	var err error

	switch v := daten.(type) {
	case string:
		// Konvertiert normalen Text in []byte
		byteData = []byte(v)
	case []uint8:
		// []uint8 ist ein Alias für []byte, also direkt zuweisen
		byteData = v
	case map[string]interface{}, []interface{}:
		// Versucht, JSON-ähnliche Strukturen in []byte zu konvertieren
		byteData, err = json.Marshal(v)
		if err != nil {
			return fmt.Errorf("fehler beim Konvertieren der Daten in JSON-Bytes: %v", err)
		}
	default:
		// Generischer Fall, falls ein anderer Typ vorliegt
		return fmt.Errorf("ungültiger Datentyp für Daten: %T", v)
	}

	// Überprüft, ob eine Seite mit dem gegebenen Titel bereits existiert.
	exists, err := CheckIfLernseiteExists(db, titel)
	if err != nil {
		return fmt.Errorf("fehler beim Überprüfen der Seite: %v", err)
	}

	var query string
	if exists {
		// Update einer existierenden Seite
		query = `
        UPDATE lernseiten 
        SET daten = $1, benutzer_id = $2
        WHERE titel = $3;
        `
		_, err = db.Exec(query, byteData, benutzerID, titel)
		if err != nil {
			return fmt.Errorf("fehler beim Aktualisieren der Seite: %v", err)
		}
		log.Printf("Seite mit dem Titel %s erfolgreich aktualisiert.\n", titel)
	} else {
		// Einfügen einer neuen Seite
		query = `
        INSERT INTO lernseiten (titel, daten, benutzer_id)
        VALUES ($1, $2, $3)
        RETURNING id;
        `
		var lastInsertID int
		err = db.QueryRow(query, titel, byteData, benutzerID).Scan(&lastInsertID)
		if err != nil {
			return fmt.Errorf("fehler beim Einfügen der neuen Seite: %v", err)
		}
		log.Printf("Neue Seite mit ID %d und Titel %s erfolgreich erstellt.\n", lastInsertID, titel)
	}

	return nil
}
