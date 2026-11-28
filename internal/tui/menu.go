package tui

import tea "charm.land/bubbletea/v2"

// entree est une ligne de menu : son libellé et la fonction qui crée l'écran à
// ouvrir. ouvrir vaut nil pour la ligne "Quitter".
type entree struct {
	libelle string
	ouvrir  func() tea.Model
}

// menuModel est un menu vertical, partagé par l'application client et
// l'application admin : seules les entrées changent.
type menuModel struct {
	titre   string
	entrees []entree
	curseur int // index de la ligne sélectionnée (zero value : la première)
}

// newMenu crée un menu avec son titre et ses entrées.
func newMenu(titre string, entrees []entree) menuModel {
	return menuModel{titre: titre, entrees: entrees}
}

// Init : rien à faire au démarrage.
func (m menuModel) Init() tea.Cmd { return nil }

// Update déplace le curseur, ouvre l'écran choisi ou quitte.
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg) // assertion de type : est-ce une touche ?
	if !ok {
		return m, nil // pas une touche : le menu n'a rien à faire
	}

	switch key.String() {
	case "q":
		return m, tea.Quit
	case "up":
		if m.curseur > 0 {
			m.curseur--
		}
	case "down":
		if m.curseur < len(m.entrees)-1 {
			m.curseur++
		}
	case "enter":
		e := m.entrees[m.curseur]
		if e.ouvrir == nil { // la ligne "Quitter"
			return m, tea.Quit
		}
		return m, Navigate(e.ouvrir()) // on crée l'écran et on demande à l'ouvrir
	}
	return m, nil
}

// View dessine le titre, les lignes (la sélectionnée en couleur) et l'aide.
func (m menuModel) View() tea.View {
	s := styleTitre.Render(m.titre) + "\n"
	for i, e := range m.entrees {
		if i == m.curseur {
			s += styleSelection.Render("> "+e.libelle) + "\n"
		} else {
			s += "  " + e.libelle + "\n"
		}
	}
	s += "\n" + styleAide.Render("↑/↓ se déplacer • Entrée valider • q quitter")
	return tea.NewView(styleCadre.Render(s))
}
