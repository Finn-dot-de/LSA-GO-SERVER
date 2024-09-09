package insert

import (
	"database/sql"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/get"
)

// Speichert oder aktualisiert eine Lernseite basierend auf der ID
func SaveLearningSitesByID(db *sql.DB, id int, titel string, daten string, benutzerID int) error {
	var query string
	var err error

	// Prüfen, ob eine Seite existiert
	exists, err := get.CheckIfLernseiteExists(db, id)
	if err != nil {
		return fmt.Errorf("Fehler beim Überprüfen der Seite: %v", err)
	}

	if exists {
		// Wenn die Seite existiert, aktualisieren
		query = `
        UPDATE lernseiten 
        SET titel = $1, daten = $2, benutzer_id = $3
        WHERE id = $4;
        `
		_, err = db.Exec(query, titel, daten, benutzerID, id)
		if err != nil {
			return fmt.Errorf("Fehler beim Aktualisieren der Seite: %v", err)
		}
		fmt.Printf("Seite mit ID %d erfolgreich aktualisiert.\n", id)
	} else {
		// Wenn die Seite nicht existiert, neue Seite einfügen
		query = `
        INSERT INTO lernseiten (titel, daten, benutzer_id)
        VALUES ($1, $2, $3)
        RETURNING id;
        `
		var lastInsertID int
		err = db.QueryRow(query, titel, daten, benutzerID).Scan(&lastInsertID)
		if err != nil {
			return fmt.Errorf("Fehler beim Einfügen der neuen Seite: %v", err)
		}
		fmt.Printf("Neue Seite mit ID %d erfolgreich erstellt.\n", lastInsertID)
	}

	return nil
}
