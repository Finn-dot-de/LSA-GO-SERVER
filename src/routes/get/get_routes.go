package get

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/get"
	"github.com/go-chi/chi"
)

// DefineGetRoutes definiert alle GET-Routen der Anwendung
func DefineGetRoutes(r *chi.Mux, db *sql.DB) {
	// Fachbezogene Fragen abrufen
	r.Get("/api/fragen/{name}", func(w http.ResponseWriter, r *http.Request) {
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
	r.Get("/api/faecher", func(w http.ResponseWriter, r *http.Request) {
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
}
