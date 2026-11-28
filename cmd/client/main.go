// Command client lance l'application client dans le terminal.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"e-commerce/internal/apiclient"
	"e-commerce/internal/tui"
)

func main() {
	// L'adresse de l'API vient de la variable d'environnement API_URL.
	url := os.Getenv("API_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	api := apiclient.New(url)

	// main ne fait que créer le client HTTP et lancer l'interface : tout le
	// reste est dans le package tui, pour pouvoir le relancer par SSH.
	if _, err := tea.NewProgram(tui.NewClientModel(api)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "erreur :", err)
		os.Exit(1)
	}
}
