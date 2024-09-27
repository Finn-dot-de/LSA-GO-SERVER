package post

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/handler_func/user_func"
	"log"
	"net/http"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/handler_func/site_func"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
	"github.com/go-chi/chi"
)

// DefinePostRoutes definiert alle POST-Routen der Anwendung
func DefinePostRoutes(r *chi.Mux, db *sql.DB) {
	r.Post("/app/api/save/lerning/site", func(w http.ResponseWriter, r *http.Request) {
		// Sicherstellen, dass der Anfrage-Body geschlossen wird, um Ressourcenlecks zu vermeiden
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Printf("Fehler beim Schließen des Anfrage-Body: %v", err)
			}
		}()

		w.Header().Set("Content-Type", "application/json")

		// Benutzer-Kürzel aus den Headern auslesen und Benutzer-ID abrufen
		kuerzel := r.Header.Get("X-Forwarded-User")
		BenutzerID, err := user_func.GetUserID(db, kuerzel)
		if err != nil {
			handleError(w, fmt.Errorf("Fehler beim Abrufen der Benutzer-ID: %v", err), http.StatusInternalServerError)
			return
		}

		// Verwende SaveRequest, um die Anfrage zu parsen
		var req structs.Lernseite
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			handleError(w, fmt.Errorf("ungültige Anfrage: %v", err), http.StatusBadRequest)
			return
		}

		// Validierung der Eingabedaten
		if req.Titel == "" {
			handleError(w, fmt.Errorf("fehlende oder ungültige Daten: Titel ist leer"), http.StatusBadRequest)
			return
		}

		// Konvertiere den Text in []byte
		textBytes := []byte(req.Text)

		// Aufruf der SaveLearningSitesByID-Funktion, um die Seite zu speichern oder zu aktualisieren
		err = site_func.SaveLearningSitesByID(db, req.Titel, textBytes, BenutzerID)
		if err != nil {
			handleError(w, fmt.Errorf("fehler beim Speichern der Seite: %v", err), http.StatusInternalServerError)
			return
		}

		// Erfolgsantwort
		sendJSONResponse(w, map[string]string{"message": "Seite erfolgreich gespeichert"})
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
