package tui

import tea "charm.land/bubbletea/v2"

// NavigateMsg demande à l'application d'afficher un autre écran.
type NavigateMsg struct {
	Screen tea.Model
}

// BackMsg demande à l'application de revenir au menu principal.
type BackMsg struct{}

// Navigate renvoie la commande qui ouvre l'écran s. Un écran l'utilise dans
// son Update : return m, Navigate(nouvelEcran).
func Navigate(s tea.Model) tea.Cmd {
	// Une tea.Cmd est une fonction qui renvoie un message : bubbletea l'exécute,
	// puis renvoie ce message à l'Update de l'application.
	return func() tea.Msg {
		return NavigateMsg{Screen: s}
	}
}

// Back est la commande qui ramène au menu principal. Comme tea.Quit, c'est
// directement une tea.Cmd : un écran écrit return m, Back.
func Back() tea.Msg {
	return BackMsg{}
}

// app est le modèle racine : il garde le menu et l'écran affiché, et fait
// suivre chaque message à l'un ou à l'autre.
type app struct {
	menu    tea.Model // le menu principal
	courant tea.Model // l'écran ouvert depuis le menu ; nil = on affiche le menu
}

// newApp crée l'application qui démarre sur le menu donné.
func newApp(menu tea.Model) app {
	return app{menu: menu}
}

// Init lance l'éventuelle commande de démarrage du menu.
func (a app) Init() tea.Cmd {
	return a.menu.Init()
}

// Update gère la navigation, puis délègue le message à l'écran visible.
func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" { // valable partout, même dans un formulaire
			return a, tea.Quit
		}
	case NavigateMsg:
		a.courant = msg.Screen
		return a, a.courant.Init() // le nouvel écran peut lancer un chargement
	case BackMsg:
		a.courant = nil
		return a, nil
	}

	// Sinon, le message va à l'écran visible, qui renvoie sa version modifiée.
	var cmd tea.Cmd
	if a.courant == nil {
		a.menu, cmd = a.menu.Update(msg)
	} else {
		a.courant, cmd = a.courant.Update(msg)
	}
	return a, cmd
}

// View affiche l'écran visible, en plein écran.
func (a app) View() tea.View {
	var v tea.View
	if a.courant == nil {
		v = a.menu.View()
	} else {
		v = a.courant.View()
	}
	v.AltScreen = true // réglé une seule fois ici, pour tous les écrans
	return v
}
