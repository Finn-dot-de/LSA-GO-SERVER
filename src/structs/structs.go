package structs

// Antwort repräsentiert eine Antwortmöglichkeit zu einer Frage.
type Antwort struct {
	AntwortID  int    `json:"antwortID"`  // ID der Antwort
	Antwort    string `json:"antwort"`    // Text der Antwort
	IstKorrekt bool   `json:"istKorrekt"` // Gibt an, ob die Antwort korrekt ist
}

// Frage repräsentiert eine Frage aus der Datenbank.
type Frage struct {
	ID          int       `json:"id"`           // ID der Frage
	FrageText   string    `json:"frageText"`    // Text der Frage
	ThemaID     int       `json:"themaID"`      // ID des Themas, zu dem die Frage gehört
	ThemaName   string    `json:"themaName"`    // Name des Themas, zu dem die Frage gehört
	Beschreibung string   `json:"beschreibung"` // Beschreibung des Themas
	Antworten   []Antwort `json:"antworten"`    // Antwortmöglichkeiten zu dieser Frage
}

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