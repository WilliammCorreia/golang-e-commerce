package tui

import "charm.land/lipgloss/v2"

// Styles communs à tous les écrans, pour que les deux applications aient le
// même aspect. Un style lipgloss est une valeur en lecture seule : Render ne le
// modifie jamais. Ces variables ne contiennent donc aucun état partagé entre
// les sessions, ce qui reste compatible avec le SSH.
var (
	// styleTitre : le titre en haut de chaque écran.
	styleTitre = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffd449")).
			MarginBottom(1)

	// styleSelection : la ligne sous le curseur.
	styleSelection = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#d70085"))

	// styleAide : les rappels de touches, en gris discret.
	styleAide = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080"))

	// styleErreur et styleSucces : les messages après une action.
	styleErreur = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	styleSucces = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))

	// styleCadre : le cadre arrondi qui entoure le contenu d'un écran.
	styleCadre = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#e800dc")).
			Padding(1, 2)
)
