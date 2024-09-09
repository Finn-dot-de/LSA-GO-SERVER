package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/Finn-dot-de/LernStoffAnwendung/src/utils"
)

func OAuth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("X-Forwarded-User")
	log.Println("Hier >>>>>>>>>>>>> ")
	if user == "" {
		log.Println("Fehler: Kein Benutzername im Header gefunden")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// JWT-Token generieren
	token, err := utils.GenerateJWT(user)
	if err != nil {
		log.Println("Fehler beim Generieren des JWT-Tokens:", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Setze den JWT-Token als HttpOnly-Cookie mit Pfad "/"
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Path:     "/",         // Der Cookie ist für alle Routen verfügbar
		Domain:   "localhost", // In der Produktion sollte hier deine tatsächliche Domain stehen
		Secure:   true,        // Setze auf true, wenn du HTTPS verwendest
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour), // Der Cookie läuft nach 24 Stunden ab
	}
	http.SetCookie(w, cookie)

	// time.Sleep(6 * time.Second)
	// Logge den Cookie, um zu überprüfen, ob er korrekt gesetzt wurde
	log.Printf("JWT Cookie gesetzt: %+v\n", cookie)

	// Leite den Benutzer nach erfolgreicher Authentifizierung weiter
	http.Redirect(w, r, "/app/", http.StatusSeeOther)
}
