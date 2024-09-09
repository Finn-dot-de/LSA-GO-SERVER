package post

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/insert"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/structs"
	"github.com/go-chi/chi"
	"net/http"
)

// DefinePostRoutes definiert alle POST-Routen der Anwendung
func DefinePostRoutes(r *chi.Mux, db *sql.DB) {
	//r.Post("/app/api/save/own/files", insert.SaveFiles) { })

	r.Options("/app/api/save/lerning/site", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080") // Anpassen auf dein Setup
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/app/api/save/lerning/site", func(w http.ResponseWriter, r *http.Request) {
		// Daten aus der Anfrage lesen
		var req structs.SaveRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
			return
		}

		// Daten an die Funktion SaveLearningSitesByID übergeben
		err = insert.SaveLearningSitesByID(db, req.ID, req.Titel, req.DateiPfad, req.BenutzerID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Fehler beim Speichern der Seite: %v", err), http.StatusInternalServerError)
			return
		}

		// Erfolgsantwort
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Seite erfolgreich gespeichert")
	})
}
