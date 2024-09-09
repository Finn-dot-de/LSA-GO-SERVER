package get

import (
	"database/sql"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
)

// GetSitesFromDB Seiten aus der Datenbank abrufen
func GetSitesFromDB(db *sql.DB) ([]structs.Lernseite, error) {
	query := `
    SELECT id, titel, datei_pfad, benutzer_id, erstellungsdatum
    FROM lernseiten;
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Abrufen der Seiten: %v", err)
	}
	defer rows.Close()

	var lernseiten []structs.Lernseite

	for rows.Next() {
		var lernseite structs.Lernseite
		err := rows.Scan(&lernseite.ID, &lernseite.Titel, &lernseite.DateiPfad, &lernseite.BenutzerID, &lernseite.Erstellungsdatum)
		if err != nil {
			return nil, fmt.Errorf("Fehler beim Scannen der Daten: %v", err)
		}
		lernseiten = append(lernseiten, lernseite)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Fehler bei der Verarbeitung der Daten: %v", err)
	}

	return lernseiten, nil
}
