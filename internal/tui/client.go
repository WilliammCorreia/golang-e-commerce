package tui

import (
	tea "charm.land/bubbletea/v2"
)

// clientModel est l'état de l'application client : les lignes du menu et la
// position du curseur. Il reste privé : cmd/client ne voit qu'un tea.Model.
type clientModel struct {
	choix   []string // les lignes du menu
	curseur int      // index de la ligne sélectionnée (zero value : 0, la première)
}

// Init est appelée une seule fois au lancement. Rien à faire ici : on renvoie nil.
func (m clientModel) Init() tea.Cmd { return nil }

// Update reçoit chaque événement (touche, redimensionnement...) et renvoie le
// modèle modifié. Récepteur valeur : on modifie une copie, que bubbletea garde.
func (m clientModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Switch de type : on regarde le vrai type du message reçu.
	switch msg := msg.(type) {
	case tea.KeyPressMsg: // c'est une touche du clavier
		// Switch classique sur le nom de la touche.
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit // on demande à bubbletea de s'arrêter
		case "up":
			if m.curseur > 0 { // on ne remonte pas au-dessus de la première ligne
				m.curseur--
			}
		case "down":
			if m.curseur < len(m.choix)-1 { // ni en dessous de la dernière
				m.curseur++
			}
		case "enter":
			if m.choix[m.curseur] == "Quitter" { // la ligne sous le curseur
				return m, tea.Quit
			}
		}
	}
	return m, nil // le modèle (peut-être modifié), et aucune commande à lancer
}

// View dessine l'écran à partir du modèle, et uniquement de lui.
func (m clientModel) View() tea.View {
	s := "Menu\n\n"
	for i, c := range m.choix {
		if i == m.curseur {
			s += "> " + c + "\n" // la ligne sélectionnée a un curseur
		} else {
			s += "  " + c + "\n"
		}
	}
	s += "\n↑/↓ pour se déplacer, Entrée pour valider, q pour quitter.\n"

	v := tea.NewView(s)
	v.AltScreen = true // plein écran : le terminal est restauré à la sortie
	return v
}

// NewClientModel crée l'interface de l'application client, avec son menu
// principal. Elle renvoie un tea.Model : l'appelant n'a pas besoin de connaître
// clientModel, seulement ses trois méthodes Init, Update et View.
func NewClientModel() tea.Model {
	return clientModel{
		choix: []string{
			"Rechercher un produit",
			"Mon panier",
			"Mes commandes",
			"Quitter",
		},
	}
}
