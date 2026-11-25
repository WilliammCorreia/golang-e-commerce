package mailer

import "log"

// LogMailer est une implémentation de Mailer qui logge les emails au lieu de les envoyer réellement.
type LogMailer struct{}

// Send écrit le message dans les logs au lieu de l'envoyer.
func (l LogMailer) Send(to, subject, body string) error {
	log.Printf("mail pour %s, objet : %s, corps : %s", to, subject, body)
	return nil
}
