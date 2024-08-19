package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/SQL"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/login"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/utils"
	"github.com/go-chi/chi"
)

func main() {
	// Verbindung zur Datenbank herstellen
	db, err := SQL.ConnectToDB()
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	r := chi.NewRouter()

	r.Use(utils.LoggerMiddleware)
	r.Use(utils.NoCacheMiddleware) // Fügen Sie die NoCacheMiddleware hinzu

	// Statische Dateien servieren (z. B. für Angular-Anwendung)
	fs := http.FileServer(http.Dir("./project"))
	r.Handle("/*", http.StripPrefix("/", fs))

	// API-Endpunkte
	r.Get("/api/fragen/{name}", func(w http.ResponseWriter, r *http.Request) {
		// Fachname aus dem Pfadparameter abrufen
		fachName := chi.URLParam(r, "name")

		// Fach aus der Datenbank abrufen
		fach, err := SQL.GetFragenFromDBNachFach(db, fachName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// JSON als Antwort senden
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(fach)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Get("/api/faecher", func(w http.ResponseWriter, r *http.Request) {

		faecher, err := SQL.GetFeacherFromDB(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(faecher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Login-Handler
	r.Post("/api/login", login.LoginHandler)

	// Server starten und auf Port 8080 lauschen
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Fallback auf Port 8080, wenn APP_PORT nicht gesetzt ist
	}
	log.Println("Der Server läuft auf Port " + appPort + "!")
	log.Fatal(http.ListenAndServe(":"+appPort, r))
}
