package user_func

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
)

// GetUserFromDB ruft die Benutzerinformationen basierend auf dem Benutzerkürzel ab
func GetUserFromDB(userkuerzel string, db *sql.DB) (*structs.Benutzer, error) {
	// Eingabevalidierung: Überprüfen, ob das Benutzerkürzel leer ist
	if userkuerzel == "" {
		return nil, errors.New("Benutzerkürzel darf nicht leer sein")
	}

	// SQL-Abfrage vorbereiten
	query := `
		SELECT benutzer.id, vorname, nachname, rolle, userkuerzel 
		FROM benutzer 
		INNER JOIN public.rollen r ON r.id = benutzer.rollen_id 
		WHERE userkuerzel = $1;
	`

	// Abfrage ausführen
	row := db.QueryRow(query, userkuerzel)

	// Benutzer-Informationen initialisieren
	var benutzer structs.Benutzer
	err := row.Scan(&benutzer.ID, &benutzer.Vorname, &benutzer.Nachname, &benutzer.Rolle, &benutzer.Userkuerzel)
	if err != nil {
		// Detaillierte Fehlerbehandlung: Überprüfen, ob keine Zeilen gefunden wurden
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Benutzer '%s' nicht gefunden", userkuerzel)
			return nil, fmt.Errorf("benutzername '%s' nicht gefunden", userkuerzel)
		}

		// Allgemeine Fehlerbehandlung und Logging
		log.Printf("Fehler beim Abrufen der Benutzerdaten für '%s': %v", userkuerzel, err)
		return nil, fmt.Errorf("fehler beim Abrufen der Benutzerdaten: %v", err)
	}

	// Erfolgreich gefundene Benutzerdaten zurückgeben
	return &benutzer, nil
}

func GetUserID(db *sql.DB, kuerzel string) (int, error) {
	var query string

	query = `SELECT id FROM benutzer WHERE userkuerzel = $1`

	var userID int
	err := db.QueryRow(query, kuerzel).Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
