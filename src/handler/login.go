package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/sql/get"
	"github.com/Finn-dot-de/LernStoffAnwendung/src/structs/login"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginData login.LoginData
	fmt.Println("Login")
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		fmt.Println("Fehler beim Dekodieren der Anmeldedaten:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	pwd, err := get.GetUserAndPasswordByUsername(loginData.Username)
	if err != nil {
		fmt.Println("Fehler beim Abrufen des Benutzers aus der Datenbank:", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwd.Password), []byte(loginData.Password))
	if err != nil {
		fmt.Println("Fehler beim Vergleichen des Passworts:", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
