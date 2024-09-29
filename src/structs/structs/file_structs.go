package structs

type SavedRequestFile struct {
	ID         int    `json:"id"`         // Seite ID (0 für neue Seite)
	Titel      string `json:"titel"`      // Titel der Seite
	DateiPfad  string `json:"dateiPfad"`  // Pfad zur Datei
	BenutzerID int    `json:"benutzerId"` // Benutzer ID, der die Seite speichert
}

type Lernseite struct {
	ID               int    `json:"id"`
	Titel            string `json:"titel"`
	Text             string `json:"text"`
	BenutzerID       int    `json:"benutzerId"`
	Erstellungsdatum string `json:"erstellungsdatum"`
}
