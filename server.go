package main

import (
	"database/sql"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/auth"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/handler/get"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/handler/post"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/middleware"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/connection"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialisiere das Logging und leite die Ausgaben sowohl in die Konsole als auch in eine Datei um
	err := middleware.InitializeLogger("app.log")
	if err != nil {
		log.Fatalf("Fehler beim Initialisieren des Loggings: %v", err)
	}

	// Verbindung zur Datenbank herstellen
	db, err := connection.ConnectToDB()
	if err != nil {
		log.Fatalf("Fehler beim Verbinden mit der Datenbank: %v", err)
	}

	// Schließt die Datenbankverbindung bei Programmende
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Fehler beim Schließen der DB:", err)
		}
	}(db)

	// Erstelle den Router
	r := chi.NewRouter()

	// Registrierung der allgemeinen Middleware-Funktionen
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.NoCacheMiddleware)

	// Route für den OAuth2-Callback (diese Route setzt den JWT-Cookie)
	r.Get("/auth", auth.OAuth2CallbackHandler) // Diese Funktion setzt den JWT-Cookie

	// Gruppe der Routen, auf die die JWT-Middleware angewendet wird
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware) // JWT-Middleware wird auf diese Routen angewendet

		// Typ-Assertion, um chi.Router in *chi.Mux umzuwandeln
		get.DefineGetRoutes(r.(*chi.Mux), db)
		post.DefinePostRoutes(r.(*chi.Mux), db)

		// Statische Dateien servieren (z. B. für Angular-Anwendung)
		fs := http.FileServer(http.Dir("C:\\DEV\\LS-ANG\\project"))
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
