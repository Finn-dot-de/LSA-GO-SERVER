package get

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/handler_func/site_func"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/handler_func/user_func"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

// DefineGetRoutes definiert alle GET-Routen der Anwendung
func DefineGetRoutes(r *chi.Mux, db *sql.DB) {
	// Fachbezogene Fragen abrufen
	r.Get("/app/api/fragen/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fachName := chi.URLParam(r, "name")

		fach, err := site_func.GetFragenFromDBByFach(db, fachName)
		if err != nil {
			handleError(w, fmt.Errorf("Fehler beim Abrufen der Fragen für Fach %s: %v", fachName, err), http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, fach)
	})

	// Alle Fächer abrufen
	r.Get("/app/api/faecher", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		faecher, err := site_func.GetFeacherFromDB(db)
		if err != nil {
			handleError(w, fmt.Errorf("Fehler beim Abrufen der Fächer: %v", err), http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, faecher)
	})

	// Route zum Abrufen des Links
	r.Get("/logout/link", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		link := struct {
			URL string `json:"url"`
		}{
			URL: "/oauth2/sign_out?rd=https%3A%2F%2Fgithub.com%2Flogout",
		}
		sendJSONResponse(w, link)
	})

	// API-Endpunkt für das Abrufen einer Datei basierend auf ihrem Titel
	r.Get("/app/api/getlernsite", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		querytitel := r.URL.Query().Get("titel")
		if querytitel == "" {
			handleError(w, errors.New("titel fehlt"), http.StatusBadRequest)
			return
		}

		lernseite, err := site_func.GetLernseiteByID(db, querytitel)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				handleError(w, errors.New("Seite nicht gefunden, bitte eine neue erstellen"), http.StatusNotFound)
				return
			}
			handleError(w, fmt.Errorf("Fehler beim Abrufen der Datei: %v", err), http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, lernseite)
	})

	// Benutzerinformationen abrufen
	r.Get("/app/api/get/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userkuerzel := r.Header.Get("X-Forwarded-User")
		userkuerzel = "Finn-dot-de"
		if userkuerzel == "" {
			handleError(w, errors.New("Kein angemeldeter Benutzer"), http.StatusBadRequest)
			return
		}

		userdata, err := user_func.GetUserFromDB(userkuerzel, db)
		if err != nil {
			if err.Error() == "Benutzername nicht gefunden" {
				handleError(w, errors.New("Benutzer nicht gefunden"), http.StatusNotFound)
			} else {
				handleError(w, errors.New("Interner Serverfehler"), http.StatusInternalServerError)
			}
			return
		}

		sendJSONResponse(w, userdata)
	})
}

// sendJSONResponse sendet die gegebene Datenstruktur als JSON-Antwort
func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Fehler beim Kodieren der JSON-Antwort: %v", err)
		http.Error(w, "Fehler beim Kodieren der Antwort", http.StatusInternalServerError)
	}
}

// handleError behandelt Fehler, protokolliert sie und sendet eine HTTP-Fehlermeldung zurück
func handleError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("Fehler: %v", err)
	http.Error(w, err.Error(), statusCode)
}
