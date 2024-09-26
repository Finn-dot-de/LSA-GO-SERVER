package post

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/handler/site_funcs"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
	"github.com/go-chi/chi"
)

// DefinePostRoutes definiert alle POST-Routen der Anwendung
func DefinePostRoutes(r *chi.Mux, db *sql.DB) {

	r.Post("/app/api/save/lerning/site", func(w http.ResponseWriter, r *http.Request) {
		// Sicherstellen, dass der Anfrage-Body geschlossen wird, um Ressourcenlecks zu vermeiden
		defer r.Body.Close()

		// Content-Type setzen
		w.Header().Set("Content-Type", "application/json")

		// Daten aus der Anfrage lesen
		var req structs.SaveRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			handleError(w, fmt.Errorf("ungültige Anfrage: %v", err), http.StatusBadRequest)
			return
		}

		// Eingabedaten validieren (Beispiel: ID und Titel dürfen nicht leer sein)
		if req.ID == 0 || req.Titel == "" {
			handleError(w, fmt.Errorf("fehlende oder ungültige Daten"), http.StatusBadRequest)
			return
		}

		// Daten an die Funktion SaveLearningSitesByID übergeben
		err = site_funcs.SaveLearningSitesByID(db, req.ID, req.Titel, req.DateiPfad, req.BenutzerID)
		if err != nil {
			handleError(w, fmt.Errorf("Fehler beim Speichern der Seite: %v", err), http.StatusInternalServerError)
			return
		}

		// Erfolgsantwort als JSON senden
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
