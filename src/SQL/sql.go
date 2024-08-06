// Das Paket SQL enthält die Funktionen zur Interaktion mit der Datenbank.
package SQL

// Importieren der notwendigen Pakete.
import (
	"database/sql" // Paket für die Interaktion mit SQL-Datenbanken.
	"errors"       // Paket für das Handling von Fehlern.
	"fmt"          // Paket für formatierte E/A.
	"log"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs" // Paket für die Structs für die JSON
	_ "github.com/lib/pq"                                   // PostgreSQL-Treiber.
)

func GetUserByUsername(username string) (structs.Password, error) {
	db, err := ConnectToDB()
	if err != nil {
		return structs.Password{}, nil
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	var user structs.User
	var pwd structs.Password
	err = db.QueryRow(`
	SELECT 
		b.name, 
		bl.passwort 
	FROM quizschema.benutzer 
		AS b 
	    	JOIN quizschema.benutzer_login 
				AS bl ON b.id = bl.id 
	WHERE b.name = $1;`,
		username,
	).Scan(&user.Username, &pwd.Password)
	log.Println(err, user)
	if err != nil {
		if err == sql.ErrNoRows {
			return structs.Password{}, errors.New("user not found")
		}
		return structs.Password{}, err
	}

	return pwd, nil
}

// ConnectToDB stellt eine Verbindung zur Datenbank her und gibt diese zurück.
func ConnectToDB() (*sql.DB, error) {
	// Konstanten für die Datenbankverbindung.
	const (
		host     = "localhost" // Der Host der Datenbank.
		port     = 5432        // Der Port der Datenbank.
		user     = "learner"   // Der Benutzername für die Datenbank.
		password = "learner"   // Das Passwort für die Datenbank.
		dbname   = "quizdb"    // Der Name der Datenbank.
	)

	// Erstellen der Verbindungszeichenkette.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Versuch, eine Verbindung zur Datenbank herzustellen.
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		// Wenn ein Fehler auftritt, gibt nil und den Fehler zurück.
		return nil, err
	}

	// Versuch, die Datenbank anzupingen.
	err = db.Ping()
	if err != nil {
		// Wenn ein Fehler auftritt, gibt nil und den Fehler zurück.
		return nil, err
	}

	// Wenn die Verbindung erfolgreich hergestellt wurde, wird eine Erfolgsmeldung gedruckt.
	fmt.Println("Successfully connected!")
	// Gibt die Datenbankverbindung und nil für den Fehler zurück.
	return db, nil
}

// GetFragenFromDBNachFach ruft alle Fragen zu einem bestimmten Thema aus der Datenbank ab und gibt sie zurück.
func GetFragenFromDBNachFach(db *sql.DB, themaName string) ([]structs.Frage, error) {
	rows, err := db.Query(`
	SELECT 
	    fragen.id, 
	    fragen.frage_text, 
	    themen.id, 
	    themen.thema_name, 
	    themen.beschreibung, 
	    a.id, 
	    a.antwort_text, 
	    a.ist_korrekt 
	FROM quizschema.fragen 
	    JOIN quizschema.moegliche_antworten 
	        AS a ON a.frage_id = fragen.id 
	    JOIN quizschema.themen 
	        ON fragen.thema_id = themen.id 
	WHERE themen.thema_name = $1;`, themaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fragen []structs.Frage
	for rows.Next() {
		var frageID, themaID int
		var frageText, themaName, beschreibung string
		var antwortID sql.NullInt64
		var antwortText sql.NullString
		var istKorrekt sql.NullBool

		// Daten aus der Abfrage in Variablen scannen
		err := rows.Scan(&frageID, &frageText, &themaID, &themaName, &beschreibung, &antwortID, &antwortText, &istKorrekt)
		if err != nil {
			return nil, err
		}

		// Überprüfen, ob die Frage bereits in der Liste vorhanden ist
		found := false
		for i := range fragen {
			if fragen[i].ID == frageID {
				found = true
				break
			}
		}

		// Wenn die Frage nicht gefunden wurde, füge sie der Liste hinzu
		if !found {
			fragen = append(fragen, structs.Frage{
				ID:           frageID,
				FrageText:    frageText,
				ThemaID:      themaID,
				ThemaName:    themaName,
				Beschreibung: beschreibung,
				Antworten:    []structs.Antwort{},
			})
		}

		// Füge die Antwort der entsprechenden Frage hinzu, falls vorhanden
		if antwortID.Valid {
			fragen[len(fragen)-1].Antworten = append(fragen[len(fragen)-1].Antworten, structs.Antwort{
				AntwortID:  int(antwortID.Int64),
				Antwort:    antwortText.String,
				IstKorrekt: istKorrekt.Bool,
			})
		}
	}
	// Überprüfen auf Fehler während des Durchlaufs der Zeilen
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fragen, nil
}
