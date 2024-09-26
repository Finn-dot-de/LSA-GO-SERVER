package site_funcs

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
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

// GetSitesFromDB Seiten aus der Datenbank abrufen
func GetSitesFromDB(db *sql.DB) ([]structs.Lernseite, error) {
	query := `
    SELECT id, titel, daten, benutzer_id, erstellungsdatum
    FROM lernseiten;
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abrufen der Seiten: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var lernseiten []structs.Lernseite

	// Hilfsfunktion für das Scannen von Rows
	err = scanLernseiten(rows, &lernseiten)
	if err != nil {
		return nil, err
	}

	return lernseiten, nil
}

// scanLernseiten scannt die Rows und fügt sie zur gegebenen Lernseiten-Liste hinzu
func scanLernseiten(rows *sql.Rows, lernseiten *[]structs.Lernseite) error {
	for rows.Next() {
		var lernseite structs.Lernseite
		err := rows.Scan(&lernseite.ID, &lernseite.Titel, &lernseite.DateiPfad, &lernseite.BenutzerID, &lernseite.Erstellungsdatum)
		if err != nil {
			return fmt.Errorf("fehler beim Scannen der Daten: %v", err)
		}
		*lernseiten = append(*lernseiten, lernseite)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("fehler bei der Verarbeitung der Daten: %v", err)
	}

	return nil
}

// GetLernseiteByID Funktion, um eine Seite basierend auf der ID aus der Datenbank abzurufen
func GetLernseiteByID(db *sql.DB, id int) (*structs.Lernseite, error) {
	query := `
		SELECT id, titel, datei_pfad, benutzer_id, erstellungsdatum
		FROM lernseiten
		WHERE id = $1;
	`

	row := db.QueryRow(query, id)

	var lernseite structs.Lernseite
	err := row.Scan(&lernseite.ID, &lernseite.Titel, &lernseite.DateiPfad, &lernseite.BenutzerID, &lernseite.Erstellungsdatum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("keine Datei mit ID %d gefunden", id)
		}
		return nil, fmt.Errorf("fehler beim Abrufen der Datei: %v", err)
	}

	return &lernseite, nil
}

// GetFragenFromDBNachFach ruft alle Fragen zu einem bestimmten Thema aus der Datenbank ab und gibt sie zurück.
func GetFragenFromDBNachFach(db *sql.DB, FachName string) ([]structs.Frage, error) {
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
func CheckIfLernseiteExists(db *sql.DB, id int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM lernseiten WHERE id=$1)"
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Überprüfen der Seite: %v", err)
	}
	return exists, nil
}

// SaveLearningSitesByID speichert oder aktualisiert eine Lernseite basierend auf der ID
func SaveLearningSitesByID(db *sql.DB, id int, titel string, daten string, benutzerID int) error {
	exists, err := CheckIfLernseiteExists(db, id)
	if err != nil {
		return fmt.Errorf("fehler beim Überprüfen der Seite: %v", err)
	}

	var query string
	if exists {
		query = `
        UPDATE lernseiten 
        SET titel = $1, daten = $2, benutzer_id = $3
        WHERE id = $4;
        `
		_, err = db.Exec(query, titel, daten, benutzerID, id)
		if err != nil {
			return fmt.Errorf("fehler beim Aktualisieren der Seite: %v", err)
		}
		log.Printf("Seite mit ID %d erfolgreich aktualisiert.\n", id)
	} else {
		query = `
        INSERT INTO lernseiten (titel, daten, benutzer_id)
        VALUES ($1, $2, $3)
        RETURNING id;
        `
		var lastInsertID int
		err = db.QueryRow(query, titel, daten, benutzerID).Scan(&lastInsertID)
		if err != nil {
			return fmt.Errorf("fehler beim Einfügen der neuen Seite: %v", err)
		}
		log.Printf("Neue Seite mit ID %d erfolgreich erstellt.\n", lastInsertID)
	}

	return nil
}
