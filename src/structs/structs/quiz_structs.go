package structs

// Antwort repräsentiert eine Antwortmöglichkeit zu einer Frage.
type Antwort struct {
	AntwortID  int    `json:"antwortID"`  // ID der Antwort
	Antwort    string `json:"antwort"`    // Text der Antwort
	IstKorrekt bool   `json:"istKorrekt"` // Gibt an, ob die Antwort korrekt ist
}

// Frage repräsentiert eine Frage aus der Datenbank.
type Frage struct {
	ID        int       `json:"id"`        // ID der Frage
	FrageText string    `json:"frageText"` // Text der Frage
	FachID    int       `json:"themaID"`   // ID des Themas, zu dem die Frage gehört
	FachName  string    `json:"themaName"` // Name des Themas, zu dem die Frage gehört
	Langform  string    `json:"langform"`  // Beschreibung des Themas
	Antworten []Antwort `json:"antworten"` // Antwortmöglichkeiten zu dieser Frage
}

type Schulfach struct {
	Schulfach string `json:"fach"`
}
