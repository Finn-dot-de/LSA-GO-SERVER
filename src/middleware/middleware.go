package middleware

import (
	"github.com/Finn-dot-de/LernStoffAnwendung/src/utils"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Definiert ein einheitliches Zeitformat als Konstante
const timeFormat = time.RFC1123

// InitializeLogger konfiguriert das Logging und leitet Ausgaben sowohl in eine Datei als auch in die Konsole um.
func InitializeLogger(logFile string) error {
	// Versucht, die Log-Datei im Append-Modus zu öffnen oder zu erstellen
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// MultiWriter erstellt eine Kombination aus Datei und Konsole
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Setzt den Log-Output auf die Kombination aus Konsole und Datei
	log.SetOutput(multiWriter)

	// Optional: Setzt das Log-Format
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Zeitstempel und Dateipfad mit Zeilennummer

	return nil
}

// LoggerMiddleware protokolliert die Details jeder eingehenden HTTP-Anfrage,
// einschließlich Startzeit, Methode, URL und Dauer der Anfrage.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Startzeit wird aufgezeichnet
		start := time.Now()
		log.Printf("Startzeit: %s | Methode: %s | URL: %s | RemoteAddr: %s",
			start.Format(timeFormat), r.Method, r.RequestURI, r.RemoteAddr)

		// defer stellt sicher, dass die Endzeit und Dauer auch dann protokolliert werden,
		// wenn der nachfolgende Handler einen Fehler verursacht oder die Funktion vorzeitig beendet wird.
		defer func() {
			end := time.Now()
			log.Printf("Endzeit: %s | Dauer: %s", end.Format(timeFormat), end.Sub(start))
		}()

		// Der nächste Handler in der Kette wird aufgerufen
		next.ServeHTTP(w, r)
	})
}

// NoCacheMiddleware verhindert das Caching von HTTP-Antworten, um sicherzustellen,
// dass der Client immer die aktuellsten Daten erhält.
func NoCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Setzt mehrere Header, um sicherzustellen, dass der Client die Antwort nicht cached.
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Der nächste Handler in der Kette wird aufgerufen
		next.ServeHTTP(w, r)
	})
}

// JWTAuthMiddleware prüft die Authentifizierung mittels JWT
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			log.Println("JWT Cookie nicht gefunden:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validiert das JWT-Token und wandelt es in den richtigen Typ um
		token, err := utils.ValidateJWT(cookie.Value)
		if err != nil || !token.Valid {
			log.Println("JWT-Token ungültig:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Versuche, die Claims als utils.CustomClaims zu interpretieren
		claims, ok := token.Claims.(*utils.CustomClaims)
		if !ok {
			log.Println("Fehler bei der Typumwandlung der Claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("JWT-Token gültig für Benutzer: %v", claims.UserID)

		// Setze die Benutzerinformationen im Context, falls nötig
		// ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		// next.ServeHTTP(w, r.WithContext(ctx))

		next.ServeHTTP(w, r)
	})
}
