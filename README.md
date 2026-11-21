# E-commerce en ligne de commande

Projet Go (M2) : une plateforme d'achat en ligne utilisable depuis le terminal.

- une **API HTTP** (`net/http`, sans framework) avec une base **SQLite** (`database/sql`) ;
- une **application client** et une **application administrateur**, deux CLI séparées construites avec Charm (bubbletea, bubbles, huh, lipgloss) ;
- un accès **SSH** aux deux applications.

## Prérequis

- Docker et Docker Compose
- Go 1.27.4 (facultatif, pour l'éditeur : tout peut se lancer dans le conteneur)

## Lancer en développement

```bash
docker compose up -d --build   # démarre le conteneur Go
docker compose exec go sh      # ouvre un shell dedans
```

Puis, dans le conteneur :

```sh
go run ./cmd/server   # API
go run ./cmd/client   # application client
go run ./cmd/admin    # application administrateur
go test ./...         # tests
```

Pour arrêter : `docker compose down`.

## Organisation du code

```
cmd/
  server/     point d'entrée de l'API
  client/     point d'entrée de l'application client
  admin/      point d'entrée de l'application administrateur
internal/
  api/        routeur, handlers, middleware d'authentification
  apiclient/  client HTTP utilisé par les deux applications
  ids/        identifiants métier (PDT-, BSK-, CMD-) et codes aléatoires
  mailer/     envoi des codes (écrits dans les logs en développement)
  model/      types métier partagés
  store/      accès à la base SQLite
  tui/        interfaces Charm
```

## Équipe

- Membre A : comptes, serveur, SSH, Docker et déploiement
- Membre B : catalogue, panier, gestion des produits et des utilisateurs
- Membre C : paiement, commandes, gestion des commandes, base des interfaces
