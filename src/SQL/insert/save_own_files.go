package insert

import (
	"database/sql"
	"fmt"
	"log"
)

// SaveFilesInDB Seite in der Datenbank speichern
func SaveFilesInDB(db *sql.DB, titel string, dateiPfad string, benutzerID int) error {
	query := `
    INSERT INTO lernseiten (titel, datei_pfad, benutzer_id)
    VALUES ($1, $2, $3)
    RETURNING id;
    `
	var lastInsertID int
	err := db.QueryRow(query, titel, dateiPfad, benutzerID).Scan(&lastInsertID)
	if err != nil {
		return fmt.Errorf("fehler beim Speichern der Seite: %v", err)
	}

	log.Printf("Seite wurde erfolgreich mit ID %d gespeichert.\n", lastInsertID)
	return nil
}
