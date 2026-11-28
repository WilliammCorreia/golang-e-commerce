package tui

import tea "charm.land/bubbletea/v2"

// bientotModel est un écran provisoire pour une fonctionnalité pas encore
// branchée. Chaque membre le remplace par son vrai écran dans client.go ou
// admin.go.
type bientotModel struct {
	titre string
}

// newBientot crée l'écran provisoire qui porte le titre donné.
func newBientot(titre string) tea.Model {
	return bientotModel{titre: titre}
}

// Init : rien à charger.
func (m bientotModel) Init() tea.Cmd { return nil }

// Update : Échap ou q ramène au menu.
func (m bientotModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "esc", "q":
			return m, Back
		}
	}
	return m, nil
}

// View affiche le titre et un message d'attente.
func (m bientotModel) View() tea.View {
	s := styleTitre.Render(m.titre) + "\n" +
		"Bientôt disponible.\n\n" +
		styleAide.Render("Échap retour au menu")
	return tea.NewView(styleCadre.Render(s))
}
