package tui

import (
	tea "charm.land/bubbletea/v2"

	"e-commerce/internal/apiclient"
)

// NewAdminModel crée l'application administrateur, sur le même principe que
// NewClientModel : un menu dont chaque entrée ouvre un écran de gestion.
func NewAdminModel(api *apiclient.Client) tea.Model {
	_ = api // utilisé dès que le premier vrai écran sera branché

	menu := newMenu("Administration", []entree{
		{"Gérer les produits", func() tea.Model { return newBientot("Gérer les produits") }},
		{"Gérer les commandes", func() tea.Model { return newBientot("Gérer les commandes") }},
		{"Gérer les utilisateurs", func() tea.Model { return newBientot("Gérer les utilisateurs") }},
		{"Quitter", nil},
	})
	return newApp(menu)
}
