package utils

import (
	"github.com/dgrijalva/jwt-go"
	"log"      // Das Paket "log" stellt Funktionen zum Schreiben von Protokollen zur Verfügung.
	"net/http" // Das Paket "net/http" stellt HTTP-Client- und Server-Implementierungen zur Verfügung.
	"time"     // Das Paket "time" stellt Funktionen zum Messen und Anzeigen von Zeit zur Verfügung.
)

// LoggerMiddleware ist eine Middleware-Funktion, die Protokollinformationen für jede HTTP-Anfrage ausgibt.
func LoggerMiddleware(next http.Handler) http.Handler {
	// Ein HandlerFunc ist ein Adapter, der eine Funktion in einen http.handler umwandelt.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Die Startzeit der Anfrage wird aufgezeichnet.
		start := time.Now()
		// Protokolliert die Startzeit, Methode, URL und Remote-Adresse der Anfrage.
		log.Printf("Startzeit: %s | Methode: %s | URL: %s | RemoteAddr: %s", start.Format(time.RFC1123), r.Method, r.RequestURI, r.RemoteAddr)
		// Ruft den nächsten handler in der Kette auf.
		next.ServeHTTP(w, r)
		// Die Endzeit der Anfrage wird aufgezeichnet.
		end := time.Now()
		// Protokolliert die Endzeit und die Dauer der Anfrage.
		log.Printf("Endzeit: %s | Dauer: %s", end.Format(time.RFC1123), end.Sub(start))
	})
}

func NoCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		log.Println(err)
		log.Println(cookie)
		if err != nil {
			log.Println("JWT Cookie nicht gefunden:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// JWT-Token validieren
		token, err := ValidateJWT(cookie.Value)
		if err != nil || !token.Valid {
			log.Println("JWT-Token ungültig:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("JWT-Token gültig für Benutzer:", token.Claims.(jwt.MapClaims)["user_id"])
		next.ServeHTTP(w, r)
	})
}