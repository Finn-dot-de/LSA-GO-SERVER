// Das Paket SQL enthält die Funktionen zur Interaktion mit der Datenbank.
package SQL

// Importieren der notwendigen Pakete.
import (
	"database/sql" // Paket für die Interaktion mit SQL-Datenbanken.
	"errors"       // Paket für das Handling von Fehlern.
	"fmt"          // Paket für formatierte E/A.
	"log"          // Paket für das Loggen von Informationen.
	"os"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs" // Paket für die Structs für die JSON-Verarbeitung.
	_ "github.com/lib/pq"                                   // PostgreSQL-Treiber.
)

// GetUserAndPasswordByUsername sucht einen Benutzer in der Datenbank anhand des Benutzernamens.
func GetUserAndPasswordByUsername(username string) (structs.Password, error) {
	// Verbindet sich mit der Datenbank.
	db, err := ConnectToDB()
	if err != nil {
		return structs.Password{}, err
	}
	// Schließt die Datenbankverbindung am Ende der Funktion.
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Error closing the database:", err)
		}
	}(db)

	// Definiert Variablen für den Benutzer und das Passwort.
	var user structs.User
	var pwd structs.Password

	// Führt die SQL-Abfrage aus, um den Benutzer und das Passwort zu finden.
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
		if err == sql.ErrNoRows {
			return structs.Password{}, errors.New("user not found")
		}
		return structs.Password{}, err
	}

	// Gibt das gefundene Passwort zurück.
	return pwd, nil
}

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

// GetFeacherFromDB ruft die Fächer aus der Datenbank ab und gibt sie als Slice zurück.
func GetFeacherFromDB(db *sql.DB) ([]structs.Schulfach, error) {
	// Führt eine SQL-Abfrage aus, um die Fächer zu erhalten.
	rows, err := db.Query(`
		SELECT
			fach
		FROM feacher;
	`)
	if err != nil {
		return nil, err
	}
	// Schließt die Rows am Ende der Funktion.
	defer rows.Close()

	var feacher []structs.Schulfach
	// Iteriert über die erhaltenen Rows.
	for rows.Next() {
		var fach string
		if err := rows.Scan(&fach); err != nil {
			return nil, err
		}

		// Überprüft, ob das Fach bereits in der Liste enthalten ist.
		found := false
		for i := range feacher {
			if feacher[i].Schulfach == fach {
				found = true
				break
			}
		}

		// Wenn das Fach nicht gefunden wurde, füge es zur Liste hinzu.
		if !found {
			feacher = append(feacher, structs.Schulfach{
				Schulfach: fach,
			})
		}
	}

	// Überprüft auf Fehler während des Durchlaufs der Zeilen.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Gibt die Liste der Fächer zurück.
	return feacher, nil
}

// GetFragenFromDBNachFach ruft alle Fragen zu einem bestimmten Thema aus der Datenbank ab und gibt sie zurück.
func GetFragenFromDBNachFach(db *sql.DB, FachName string) ([]structs.Frage, error) {
	// Führt eine SQL-Abfrage aus, um die Fragen zu einem bestimmten Fach zu erhalten.
	rows, err := db.Query(`
	SELECT
	    fragen.id,
	    fragen.frage_text,
	    feacher.id,
	    feacher.fach,
	    feacher.beschreibung,
	    a.id,
	    a.antwort_text,
	    a.ist_korrekt
	FROM fragen
	    JOIN moegliche_antworten AS a ON a.frage_id = fragen.id
	    JOIN feacher ON fragen.fach_id = feacher.id
	WHERE feacher.fach = $1;`, FachName)
	if err != nil {
		return nil, err
	}
	// Schließt die Rows am Ende der Funktion.
	defer rows.Close()

	var fragen []structs.Frage
	// Iteriert über die erhaltenen Rows.
	for rows.Next() {
		var frageID, themaID int
		var frageText, themaName, beschreibung string
		var antwortID sql.NullInt64
		var antwortText sql.NullString
		var istKorrekt sql.NullBool

		// Daten aus der Abfrage in Variablen scannen.
		err := rows.Scan(&frageID, &frageText, &themaID, &themaName, &beschreibung, &antwortID, &antwortText, &istKorrekt)
		if err != nil {
			return nil, err
		}

		// Überprüfen, ob die Frage bereits in der Liste vorhanden ist.
		found := false
		for i := range fragen {
			if fragen[i].ID == frageID {
				found = true
				break
			}
		}

		// Wenn die Frage nicht gefunden wurde, füge sie der Liste hinzu.
		if !found {
			fragen = append(fragen, structs.Frage{
				ID:           frageID,
				FrageText:    frageText,
				FachID:       themaID,
				FachName:     themaName,
				Beschreibung: beschreibung,
				Antworten:    []structs.Antwort{},
			})
		}

		// Füge die Antwort der entsprechenden Frage hinzu, falls vorhanden.
		if antwortID.Valid {
			fragen[len(fragen)-1].Antworten = append(fragen[len(fragen)-1].Antworten, structs.Antwort{
				AntwortID:  int(antwortID.Int64),
				Antwort:    antwortText.String,
				IstKorrekt: istKorrekt.Bool,
			})
		}
	}
	// Überprüfen auf Fehler während des Durchlaufs der Zeilen.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Gibt die Liste der Fragen zurück.
	return fragen, nil
}
