package login

// LoginData repräsentiert die Anmeldedaten eines Benutzers.
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User repräsentiert einen Benutzer in der Datenbank.
type User struct {
	Username string
}

type Password struct {
	Password string
}
