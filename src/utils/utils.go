package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

// JWT-Keys aus Umgebungsvariablen laden, falls nicht gesetzt, wird ein Default-Wert genutzt (Achtung: Für Produktion anpassen)
var jwtKey = []byte(getJWTSecret())

// Ablaufzeit des Tokens als Konstante
const tokenExpiryDuration = 24 * time.Hour

// CustomClaims definiert die Struktur der Claims im JWT
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT generiert ein JWT mit benutzerdefinierten Claims
func GenerateJWT(userID string) (string, error) {
	// Claims definieren, einschließlich Ablaufzeit
	claims := &CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpiryDuration).Unix(),
		},
	}

	// Token mit den Claims erstellen
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Token signieren
	return token.SignedString(jwtKey)
}

// ValidateJWT validiert das gegebene JWT-Token und gibt es zurück, falls es gültig ist
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// Token parsen und die Signaturmethode prüfen
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Prüfen, ob das Token die richtige Signaturmethode verwendet
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unerwartete Signaturmethode")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Token Claims validieren
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return nil, errors.New("Token ist abgelaufen")
		}
		return token, nil
	}

	return nil, errors.New("ungültiges Token")
}

// getJWTSecret lädt den JWT-Schlüssel aus Umgebungsvariablen
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Warnung: In der Produktion sollte immer ein sicherer Schlüssel verwendet werden
		secret = "dein_geheimes_key" // Standardwert, falls Umgebungsvariable nicht gesetzt ist
	}
	return secret
}
