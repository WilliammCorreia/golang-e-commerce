# Instructions pour l'IA : projet Go e-commerce

Ce fichier s'adresse à l'assistant IA d'un membre de l'équipe. Claude Code le lit automatiquement dans ce dépôt. Pour une IA en ligne (claude.ai par exemple), colle-le au début de la conversation, avec `Sujet.pdf` et `Cours.pdf` en pièces jointes : ces deux documents font foi.

## Contexte

Projet de groupe de M2 : une plateforme d'achat en ligne utilisable depuis le terminal. Les trois étudiants débutent en Go. Le rendu est une archive ZIP, à déposer le lundi 30 novembre 2026. Il est suivi d'une soutenance où chacun doit expliquer son propre code.

| Membre | Périmètre |
|---|---|
| William (WilliammCorreia) | comptes (inscription, confirmation, connexion, mot de passe oublié), serveur HTTP, gestion des utilisateurs côté admin, SSH, Docker |
| Vincent (VincentG-C) | catalogue et recherche, panier, gestion des produits côté admin |
| Yasmine (yasyasbdr) | paiement, commandes, gestion des commandes côté admin, base des interfaces Charm |

Si tu ne sais pas avec qui tu travailles, demande-le au début. Reste dans son périmètre : ne modifie pas le code d'un autre domaine. Si c'est nécessaire, dis à l'étudiant de prévenir d'abord le propriétaire de ce code.

## Mode de travail

Le sujet pénalise le code non maîtrisé et l'usage abusif de l'IA. L'étudiant doit pouvoir expliquer chaque ligne à l'oral.

- Réponds en français.
- Avance par petites étapes. Pour chacune :
  1. explique la notion en citant le module du cours qui en parle ;
  2. donne un exemple court et commenté, de 30 lignes au plus ;
  3. laisse l'étudiant l'adapter au projet ;
  4. relis son code.
- Ne génère jamais une fonctionnalité entière d'un coup. Si l'étudiant demande du code complet, découpe-le et explique chaque morceau.
- Après chaque étape, pose une ou deux questions pour vérifier qu'il a compris.
- Reste dans ce que le cours enseigne : bases, structs, interfaces, erreurs, packages, tests, `net/http`. Si tu as besoin d'une notion qui n'y figure pas, explique-la avant de t'en servir.
- Avant d'écrire du code, fais faire un `git pull` et lis le code existant. Réutilise ce qui existe, et ne crée jamais une deuxième version de quelque chose qui existe déjà.

## Contraintes du sujet (non négociables)

- **API** : uniquement `net/http`, sans framework. Utilise les motifs du `ServeMux` de Go 1.22+ : `"GET /products/{id}"`, puis `r.PathValue("id")`.
- **Base de données** : `database/sql` avec SQLite, via `modernc.org/sqlite` (pilote `"sqlite"`, sans CGO).
  - Pas d'ORM et pas d'outil de migration.
  - Le schéma est créé au démarrage avec `CREATE TABLE IF NOT EXISTS`.
- **Interfaces** : deux applications en ligne de commande séparées, `cmd/client` et `cmd/admin`.
  - Elles sont construites avec bubbletea, bubbles, huh et lipgloss, et les quatre doivent servir.
  - Elles passent **toujours** par l'API HTTP, jamais directement par la base.
- **SSH** : les deux applications doivent aussi être accessibles par SSH, avec wish. Les interfaces doivent donc pouvoir se lancer depuis une fonction (voir la section Charm).
- **Dépendances** : aucune autre que `charm.land/*` et `modernc.org/sqlite`. Tout le reste se fait avec la bibliothèque standard. Sont notamment interdits :
  - les routeurs et frameworks web : gin, echo, chi, gorilla/mux, fiber ;
  - les couches SQL et les outils de migration : GORM, sqlx, sqlc, goose, golang-migrate ;
  - les bibliothèques de ligne de commande et de configuration : cobra, viper, godotenv ;
  - testify, google/uuid, les bibliothèques JWT, et bcrypt de golang.org/x/crypto.

## Organisation du code

```
cmd/server          point d'entrée de l'API
cmd/client          application client
cmd/admin           application administrateur
internal/api        routeur, handlers, middleware, helpers JSON
internal/apiclient  client HTTP utilisé par les deux applications
internal/ids        identifiants métier et codes aléatoires
internal/mailer     envoi des codes (écrits dans les logs)
internal/model      types partagés et leurs tags JSON
internal/password   hachage et vérification des mots de passe
internal/store      schéma SQL et requêtes
internal/tui        interfaces Charm et styles communs
```

- **Un dossier = un package.** Chaque partie du code a sa place :
  - une route API est un handler dans `internal/api/<domaine>.go`, enregistré dans `NewRouter` (`internal/api/router.go`) ;
  - les requêtes SQL vont dans `internal/store/<domaine>.go` ;
  - les types partagés, avec leurs tags JSON, vont dans `internal/model`.
- **Chaque membre écrit toute la chaîne de ses fonctionnalités** :
  1. la table et les requêtes, dans `store` ;
  2. la route, dans `api` ;
  3. la méthode du client HTTP, dans `apiclient` ;
  4. l'écran, dans `tui`.
- **Les fondations communes ont un seul propriétaire.** Ce sont la connexion à la base, la struct qui porte les dépendances des handlers, le middleware d'authentification et le client HTTP de base. S'il en manque une pour avancer, ne l'invente pas : dis à l'étudiant d'en parler à l'équipe.

## Ce qui existe déjà (à réutiliser)

- **`internal/api/json.go`** : `writeJSON(w, status, v)` et `writeError(w, status, msg)`. `writeError` renvoie `{"error": "msg"}`. Toutes les réponses de l'API passent par ces deux fonctions.
- **`internal/api/router.go`** : `NewRouter() http.Handler`, avec la route `GET /health`.
- **`internal/password`** : `Hash(password) (string, error)` et `Check(password, stored) bool`, basées sur PBKDF2. Ne hache jamais un mot de passe autrement.
- **`cmd/server/main.go`** : `getenv(key, fallback)` et un `http.Server` avec des délais d'expiration.

## Conventions de l'API

- **Format** : JSON partout, avec des noms de champs en snake_case (`price_ht`, `created_at`).
- **Codes HTTP** :
  - 200 pour une lecture ou une modification, 201 pour une création, 204 pour une suppression ;
  - 400 pour une donnée invalide, 401 si l'utilisateur n'est pas connecté, 403 s'il n'est pas administrateur ;
  - 404 si la ressource est introuvable, 409 en cas de conflit (un email déjà utilisé, par exemple).
- **Authentification** : un jeton de session opaque, envoyé dans l'en-tête `Authorization: Bearer <jeton>`.
  - Le middleware est écrit par William, et c'est lui qui fournit au handler l'utilisateur connecté. Ne revérifie pas le jeton dans tes handlers.
  - Les routes d'administration exigent un compte administrateur.
- **Données sensibles** : aucune réponse ne doit contenir un mot de passe haché, le jeton d'un autre utilisateur, un numéro de carte ou un CVC.

## Règles métier

- **Argent** : toujours en `int64` de centimes, jamais en `float64`.
  - Chaque produit a un prix HT et un taux de TVA en points de base : `2000` = 20 %, `550` = 5,5 %.
  - Le prix TTC se calcule à un seul endroit, une fonction de `internal/model`, avec la formule `(ht*(10000+taux) + 5000) / 10000`. Le `+ 5000` arrondit au centime le plus proche.
  - Un total vaut la somme des (prix TTC unitaire × quantité).
- **Identifiants métier** : un préfixe, un tiret, puis 6 caractères parmi `0-9A-Z`, tirés avec `crypto/rand`.
  - `PDT-7D2K8N` pour un produit, `BSK-1KH8E7` pour un panier, `CMD-1F2S8B` pour une commande.
  - Une seule implémentation, dans `internal/ids`. La colonne est `UNIQUE`.
  - C'est cet identifiant qu'on affiche et qu'on met dans les URL : `GET /products/PDT-7D2K8N`.
- **Recherche** de produits, sur cinq critères qui peuvent se combiner :
  - le nom et la description : recherche partielle, sans tenir compte des majuscules ;
  - la catégorie ;
  - le prix HT et le prix TTC : une fourchette minimum/maximum.
- **Panier** : un panier par utilisateur, lié à son compte et enregistré en base.
  - On peut y ajouter un ou plusieurs produits, et modifier ou supprimer une quantité.
  - On affiche le total TTC. La livraison est toujours gratuite.
- **Paiement** : on paie un panier enregistré, avec une validation locale (pas de Stripe).
  - Le numéro de carte doit passer l'algorithme de Luhn, la date d'expiration (MM/AA) doit être dans le futur, et le CVC doit faire 3 chiffres.
  - Ne stocke jamais le numéro complet ni le CVC : au plus les 4 derniers chiffres.
  - Un paiement réussi crée la commande au statut `paid` et vide le panier.
- **Commandes** : une commande reprend les produits du panier, avec la date et un identifiant `CMD-`.
  - Le prix unitaire est figé au moment de l'achat.
  - Les statuts sont des constantes de `internal/model`, en anglais dans la base et le JSON, en français à l'écran :

    | Valeur | Affichage |
    |---|---|
    | `pending` | en attente |
    | `paid` | payée |
    | `shipping` | en cours de livraison |
    | `delivered` | livrée |
    | `cancelled` | annulée, avec une raison obligatoire |
  - Une commande créée par un administrateur démarre au statut `pending`.
- **Administration** : un administrateur peut :
  - ajouter un produit, dont l'identifiant `PDT-` est généré ;
  - changer le statut d'une commande, et créer une commande pour n'importe quel client à partir d'un ou plusieurs produits ;
  - créer, modifier, supprimer et confirmer un compte utilisateur.

## SQL

- **Requêtes paramétrées uniquement**, avec `?`. Ne concatène jamais une valeur saisie dans une requête : c'est la porte ouverte aux injections SQL.
- **Lectures** : toujours `defer rows.Close()`, puis vérifier `rows.Err()` après la boucle.
- **Plusieurs écritures liées**, comme créer une commande et vider le panier : une transaction, avec `db.Begin`, `defer tx.Rollback()` et `tx.Commit()`.
- **Ligne introuvable** : `store` renvoie une erreur que le handler teste avec `errors.Is`, pour répondre 404.

## Charm v2 : attention aux exemples en ligne

Les tutoriels et la plupart des IA donnent du code v1, avec les imports `github.com/charmbracelet/...`. Ce code ne compile pas ici.

- **Imports à utiliser** (versions vérifiées le 22 novembre 2026) :

  | Import | Version |
  |---|---|
  | `tea "charm.land/bubbletea/v2"` | v2.1.2 |
  | `charm.land/bubbles/v2/<composant>` (par exemple `.../textinput`, `.../list`, `.../table`) | v2.2.2 |
  | `charm.land/huh/v2` | v2.0.3 |
  | `charm.land/lipgloss/v2` | v2.1.1 |
  | `charm.land/wish/v2` | v2.0.5 |
- **La vue** : en v2, `View()` renvoie un `tea.View` et non une `string`. On la crée avec `return tea.NewView(s)`. Le plein écran se règle sur la vue (`v.AltScreen = true`), pas avec une option du programme.
- **Les touches** arrivent sous forme de `tea.KeyPressMsg`. `msg.String()` donne `"q"`, `"ctrl+c"`, `"enter"`, etc.
- **Les formulaires huh** : `*huh.Form` ne satisfait pas `tea.Model`, car son `View()` renvoie une `string` et son `Update` renvoie `(huh.Model, tea.Cmd)`. On l'intègre donc dans son propre modèle :
  1. appeler `m.form.Update(msg)` ;
  2. récupérer le formulaire avec l'assertion de type `form.(*huh.Form)` ;
  3. afficher `m.form.View()` dans sa vue.

Pour que le SSH fonctionne :

- **Un constructeur par application.** Chaque application se construit avec une fonction qui prend le client HTTP et renvoie un `tea.Model`, par exemple `NewModel(api *apiclient.Client) tea.Model`. `cmd/client` et `cmd/admin` ne font que créer le client et lancer `tea.NewProgram(...).Run()`.
- **Tout passe par le modèle bubbletea.** Dans les interfaces, pas de `form.Run()` de huh, pas de `fmt.Println`, pas de lecture de `os.Stdin`, pas de variable globale. Sinon, l'écran s'afficherait sur le serveur au lieu de la session SSH.

## Style Go (repris du cours)

- **Erreurs** : ce sont des valeurs, enveloppées avec `fmt.Errorf("contexte : %w", err)` et testées avec `errors.Is` (module 7). Pas de `panic` pour une erreur normale.
- **Interfaces** : petites, et définies du côté de celui qui les utilise (module 6).
- **Dépendances** : passées en paramètre ou en champ de struct, jamais en variable globale.
- **Commentaires** : sur chaque nom exporté, et sur chaque package dans son `doc.go`.
- **Tests** : en table avec `t.Run` (module 9), pour la logique pure : TTC, Luhn, date d'expiration, identifiants, filtres de recherche.
- **Configuration** : uniquement par variables d'environnement : `LISTEN_ADDR`, `DB_PATH`, `API_URL`, `ADMIN_EMAIL`, `ADMIN_PASSWORD`.

## Avant chaque push

Dans le conteneur (`docker compose up -d`, puis `docker compose exec go sh`) :

```sh
go fmt ./...
go vet ./...
go test ./...
```

Ensuite, depuis ta machine :

1. Lance `git pull --rebase`, puis `git push`.
2. Fais de petits commits, avec un message en français au format `feat(domaine): Ajoute ...` ou `fix(domaine): Corrige ...`.
3. Ne commite jamais `.env` ni un fichier `.db`.
