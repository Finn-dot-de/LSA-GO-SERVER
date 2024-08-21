package get

import (
	sql2 "database/sql"
	"errors"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/connection"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/login"
	"log"
)

// GetUserAndPasswordByUsername sucht einen Benutzer in der Datenbank anhand des Benutzernamens.
func GetUserAndPasswordByUsername(username string) (login.Password, error) {
	// Verbindet sich mit der Datenbank.
	db, err := connection.ConnectToDB()
	if err != nil {
		return login.Password{}, err
	}
	// Schließt die Datenbankverbindung am Ende der Funktion.
	defer func(db *sql2.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Error closing the database:", err)
		}
	}(db)

	// Definiert Variablen für den Benutzer und das Passwort.
	var user login.User
	var pwd login.Password

	// Führt die sql-Abfrage aus, um den Benutzer und das Passwort zu finden.
	err = db.QueryRow(`
		SELECT
			b.name,
			bl.passwort
		FROM benutzer AS b
			JOIN benutzer_login AS bl ON b.id = bl.id
		WHERE b.name = $1;`,
		username,
	).Scan(&user.Username, &pwd.Password)

	// Loggt das Ergebnis der Abfrage.
	log.Println(err, user)

	// Überprüft, ob ein Fehler aufgetreten ist.
	if err != nil {
		if err == sql2.ErrNoRows {
			return login.Password{}, errors.New("user not found")
		}
		return login.Password{}, err
	}

	// Gibt das gefundene Passwort zurück.
	return pwd, nil
}
