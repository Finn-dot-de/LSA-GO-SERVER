package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/auth"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/routes/get"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/routes/post"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/connection"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/utils"
	"github.com/go-chi/chi"
)

func main() {
	// Verbindung zur Datenbank herstellen
	db, err := connection.ConnectToDB()
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	r := chi.NewRouter()

	// Registrierung der Middleware-Funktionen
	r.Use(utils.LoggerMiddleware)
	r.Use(utils.NoCacheMiddleware)
	r.Use(utils.JWTAuthMiddleware)

	// Route für den OAuth2-Callback
	r.Get("/oauth2/callback", auth.OAuth2CallbackHandler)

	// GET- und POST-Routen definieren
	get.DefineGetRoutes(r, db)
	post.DefinePostRoutes(r)

	// Statische Dateien servieren (z. B. für Angular-Anwendung)
	fs := http.FileServer(http.Dir("./project"))
	r.Handle("/*", http.StripPrefix("/", fs))

	// Server starten und auf Port 8080 lauschen
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Fallback auf Port 8080, wenn APP_PORT nicht gesetzt ist
	}
	log.Println("Der Server läuft auf Port " + appPort + "!")
	log.Fatal(http.ListenAndServe(":"+appPort, r))
}
