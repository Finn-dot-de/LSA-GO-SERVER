package structs

type SaveRequest struct {
	ID         int    `json:"id"`         // Seite ID (0 für neue Seite)
	Titel      string `json:"titel"`      // Titel der Seite
	DateiPfad  string `json:"dateiPfad"`  // Pfad zur Datei
	BenutzerID int    `json:"benutzerId"` // Benutzer ID, der die Seite speichert
}

type Lernseite struct {
	ID               int
	Titel            string
	DateiPfad        string
	BenutzerID       int
	Erstellungsdatum string
}
