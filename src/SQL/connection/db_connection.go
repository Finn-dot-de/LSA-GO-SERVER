package connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

// ConnectToDB stellt eine Verbindung zur Datenbank her und gibt diese zurück.
func ConnectToDB() (*sql.DB, error) {
	// Laden der Umgebungsvariablen.
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Erstellen der Verbindungszeichenkette.
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Versuch, eine Verbindung zur Datenbank herzustellen.
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Versuch, die Datenbank anzupingen, um die Verbindung zu testen.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Wenn die Verbindung erfolgreich hergestellt wurde, wird eine Erfolgsmeldung gedruckt.
	fmt.Println("Successfully connected!")

	// Gibt die Datenbankverbindung und nil für den Fehler zurück.
	return db, nil
}
