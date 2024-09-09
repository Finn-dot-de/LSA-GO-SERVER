package get

import (
	"database/sql"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
)

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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("keine Datei mit ID %d gefunden", id)
		}
		return nil, fmt.Errorf("Fehler beim Abrufen der Datei: %v", err)
	}

	return &lernseite, nil
}
