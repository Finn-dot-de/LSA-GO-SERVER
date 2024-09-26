package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/utils"
)

// OAuth2CallbackHandler behandelt den Callback der OAuth2-Authentifizierung.
func OAuth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("X-Forwarded-User")
	if user == "" {
		log.Println("Fehler: Kein Benutzername im Header gefunden")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// JWT-Token generieren
	token, err := utils.GenerateJWT(user)
	if err != nil {
		log.Printf("Fehler beim Generieren des JWT-Tokens für Benutzer %s: %v", user, err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Cookie-Einstellungen dynamisch anpassen (z.B. je nach Umgebungsvariable)
	secureCookie := false
	if os.Getenv("ENVIRONMENT") == "production" {
		secureCookie = true
	}

	// Setze den JWT-Token als HttpOnly-Cookie mit angepassten Einstellungen
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Domain:   getDomain(), // Funktion zur Bestimmung der Domain
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour), // Der Cookie läuft nach 24 Stunden ab
	}
	http.SetCookie(w, cookie)

	// Logge den erfolgreichen Setzvorgang des Cookies
	log.Printf("JWT Cookie für Benutzer %s gesetzt: %+v\n", user, cookie)

	// Leite den Benutzer nach erfolgreicher Authentifizierung weiter
	http.Redirect(w, r, "/app/", http.StatusSeeOther)
}

// getDomain gibt die Domain basierend auf der Umgebung zurück
func getDomain() string {
	if os.Getenv("ENVIRONMENT") == "production" {
		return "your-production-domain.com" // Ersetze dies durch deine Produktionsdomain
	}
	return "localhost"
}
