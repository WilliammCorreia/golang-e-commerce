package tui

import (
	tea "charm.land/bubbletea/v2"

	"e-commerce/internal/apiclient"
)

// NewClientModel crée l'application client : son menu principal, dont chaque
// entrée ouvre un écran. api est le client HTTP que les écrans utilisent pour
// parler au serveur ; il est capturé par les fonctions ouvrir.
//
// Pour brancher un écran, remplacer son newBientot par le vrai constructeur,
// par exemple : func() tea.Model { return newOrdersScreen(api) }.
func NewClientModel(api *apiclient.Client) tea.Model {
	_ = api // utilisé dès que le premier vrai écran sera branché

	menu := newMenu("Bienvenue dans notre boutique ♥", []entree{
		{"Rechercher un produit", func() tea.Model { return newBientot("Rechercher un produit") }},
		{"Mon panier", func() tea.Model { return newBientot("Mon panier") }},
		{"Mes commandes", func() tea.Model { return newBientot("Mes commandes") }},
		{"Quitter", nil},
	})
	return newApp(menu)
}
