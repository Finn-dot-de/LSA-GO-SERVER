package get

import (
	"database/sql"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
)

// GetUserFromDB ruft die Benutzer-ID basierend auf dem Benutzernamen ab
func GetUserFromDB(userkuerzel string, db *sql.DB) (*structs.Benutzer, error) {
	query := "SELECT benutzer.id, vorname, nachname, rolle, userkuerzel FROM benutzer inner join public.rollen r on r.id = benutzer.rollen_id WHERE userkuerzel = $1;"
	row := db.QueryRow(query, userkuerzel)

	// Benutzer-Informationen speichern
	var benutzer structs.Benutzer
	err := row.Scan(&benutzer.ID, &benutzer.Vorname, &benutzer.Nachname, &benutzer.Rolle, &benutzer.Userkuerzel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Benutzername '%s' nicht gefunden", userkuerzel)
		}
		return nil, fmt.Errorf("Fehler beim Abrufen der Benutzerdaten: %v", err)
	}

	return &benutzer, err
}
