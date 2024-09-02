package auth

import (
	"encoding/json"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/utils"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

// OAuth2CallbackHandler OAuth2 Callback Handler
func OAuth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Angenommen, der Benutzername kommt vom OAuth2-Proxy im Header
	user := r.Header.Get("X-Auth-User")
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

	// Setze den JWT-Token als HttpOnly-Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})

	// Leite den Benutzer nach erfolgreicher Authentifizierung weiter
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("JWT-Token erfolgreich gesetzt und Benutzer weitergeleitet")
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie aus der Anfrage holen
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// JWT-Token validieren
	token, err := utils.ValidateJWT(cookie.Value)
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Benutzerinformationen aus den Claims extrahieren
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	email := claims["email"].(string)

	// Benutzerinformationen als JSON zurückgeben
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": userID,
		"email":   email,
	})
}
