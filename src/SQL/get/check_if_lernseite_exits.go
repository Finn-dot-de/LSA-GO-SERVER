package get

import (
	"database/sql"
	"fmt"
)

// Funktion, um zu überprüfen, ob eine Seite existiert
func CheckIfLernseiteExists(db *sql.DB, id int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM lernseiten WHERE id=$1)"
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("Fehler beim Überprüfen der Seite: %v", err)
	}
	return exists, nil
}
