package main

import (
	"database/sql"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/auth"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/routes/get"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/routes/post"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/connection"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/utils"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
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
			log.Println("Fehler beim Schließen der DB:", err)
		}
	}(db)

	// Erstelle den Router
	r := chi.NewRouter()

	// Registrierung der allgemeinen Middleware-Funktionen
	r.Use(utils.LoggerMiddleware)
	r.Use(utils.NoCacheMiddleware)

	// Route für den OAuth2-Callback (diese Route setzt den JWT-Cookie)
	r.Get("/", auth.OAuth2CallbackHandler) // Diese Funktion setzt den JWT-Cookie

	// Gruppe der Routen, auf die die JWT-Middleware angewendet wird
	r.Group(func(r chi.Router) {
		r.Use(utils.JWTAuthMiddleware) // JWT-Middleware wird auf diese Routen angewendet

		// Typ-Assertion, um chi.Router in *chi.Mux umzuwandeln
		get.DefineGetRoutes(r.(*chi.Mux), db)
		post.DefinePostRoutes(r.(*chi.Mux), db)

		// Statische Dateien servieren (z. B. für Angular-Anwendung)
		fs := http.FileServer(http.Dir("C:\\DEV\\LS-ANG\\project"))
		print("fs: %v", fs)
		r.Handle("/*", http.StripPrefix("/app/", fs))
	})

	// Server starten und auf Port 8080 lauschen
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Fallback auf Port 8080, wenn APP_PORT nicht gesetzt ist
	}
	log.Println("Der Server läuft auf Port " + appPort + "!")
	log.Fatal(http.ListenAndServe(":"+appPort, r))
}
