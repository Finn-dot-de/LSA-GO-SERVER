package get

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/get"
	"github.com/go-chi/chi"
)

// DefineGetRoutes definiert alle GET-Routen der Anwendung
func DefineGetRoutes(r *chi.Mux, db *sql.DB) {
	// Fachbezogene Fragen abrufen
	r.Get("/app/api/fragen/{name}", func(w http.ResponseWriter, r *http.Request) {
		fachName := chi.URLParam(r, "name")

		fach, err := get.GetFragenFromDBNachFach(db, fachName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(fach)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Alle Fächer abrufen
	r.Get("/app/api/faecher", func(w http.ResponseWriter, r *http.Request) {
		faecher, err := get.GetFeacherFromDB(db)
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

	// Route zum Abrufen des Links
	r.Get("/logout/link", func(w http.ResponseWriter, r *http.Request) {
		link := struct {
			URL string `json:"url"`
		}{
			URL: "/oauth2/sign_out?rd=https%3A%2F%2Fgithub.com%2Flogout", // Dynamisch generierter oder statischer Link
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(link)
		if err != nil {
			return
		}
	})

	// API-Endpunkt für das Abrufen einer Datei basierend auf ihrer ID
	r.Get("/app/api/getlernsite", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "ID fehlt", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Ungültige ID", http.StatusBadRequest)
			return
		}

		// Datei aus der Datenbank abrufen
		lernseite, err := get.GetLernseiteByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Keine Seite gefunden, neue Seite erstellen
				http.Error(w, "Seite nicht gefunden, bitte eine neue erstellen", http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprintf("Fehler beim Abrufen der Datei: %v", err), http.StatusInternalServerError)
			return
		}

		// Dateiinformationen als JSON zurückgeben
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lernseite)
	})

	r.Get("/app/api/get/user", func(w http.ResponseWriter, r *http.Request) {
		// Benutzer-Kürzel aus dem Header auslesen
		userkuerzel := r.Header.Get("X-Forwarded-User")

		// Prüfen, ob der Benutzer-Kürzel leer ist
		if userkuerzel == "" {
			http.Error(w, "Kein angemeldeter Benutzer", http.StatusBadRequest)
			return
		}

		// Benutzerinformationen aus der Datenbank abrufen
		userdata, err := get.GetUserFromDB(userkuerzel, db)
		log.Println(err)
		if err != nil {
			if err.Error() == "Benutzername nicht gefunden" {
				http.Error(w, "Benutzer nicht gefunden", http.StatusNotFound)
			} else {
				http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
			}
			return
		}

		// Benutzerinformationen als JSON zurückgeben
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userdata); err != nil {
			http.Error(w, "Fehler beim Kodieren der Antwort", http.StatusInternalServerError)
		}
	})

}
