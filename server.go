package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	defer db.Close()

	r := chi.NewRouter()

	r.Use(utils.LoggerMiddleware)

	// Statische Dateien servieren (z. B. für Angular-Anwendung)
	fs := http.FileServer(http.Dir("./project"))
	r.Handle("/*", http.StripPrefix("/", fs))

	// API-Endpunkte
	r.Get("/api/fragen", func(w http.ResponseWriter, r *http.Request) {
		// Fragen aus der Datenbank abrufen
		fragen, err := SQL.GetFragenFromDB(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// JSON als Antwort senden
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(fragen)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Login-Handler
	r.Post("/api/login", login.LoginHandler)

	// Server starten und auf Port 8080 lauschen
	log.Println("Der Server läuft auf 8080!!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
