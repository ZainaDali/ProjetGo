# URL Shortener - Service de Raccourcissement d'URLs

## 🎯 Description

Service web performant de raccourcissement et de gestion d'URLs développé en Go. L'application permet de transformer une URL longue en une URL courte et unique, avec redirection instantanée et analytics asynchrones. Le système inclut également un moniteur pour vérifier périodiquement la disponibilité des URLs.

## ✨ Fonctionnalités

### API REST
- ✅ **GET /health** : Vérification de l'état de santé du service
- ✅ **POST /api/v1/links** : Création d'une nouvelle URL courte
- ✅ **GET /{shortCode}** : Redirection vers l'URL originale (HTTP 302)
- ✅ **GET /api/v1/links/{shortCode}/stats** : Statistiques d'un lien (nombre de clics)

### Interface CLI
- ✅ **create** : Création d'une URL courte depuis la ligne de commande
- ✅ **stats** : Affichage des statistiques d'un lien
- ✅ **migrate** : Exécution des migrations de base de données
- ✅ **run-server** : Lancement du serveur API avec workers et moniteur

### Caractéristiques Techniques
- 🔄 **Analytics asynchrones** : Enregistrement des clics en arrière-plan sans bloquer la redirection
- 📊 **Monitoring d'URLs** : Vérification périodique de la disponibilité des URLs
- 🎲 **Génération de codes uniques** : Codes courts de 6 caractères alphanumériques
- 💾 **Persistance SQLite** : Base de données légère avec GORM
- ⚙️ **Configuration flexible** : Gestion via fichier YAML et Viper

## 🚀 Installation et Démarrage

### Prérequis
- Go 1.24.3 ou supérieur
- Git

### 1. Installation

```bash
# Cloner le projet
git clone https://github.com/axellelanca/urlshortener.git
cd ProjetGo

# Télécharger les dépendances

```bash
go mod tidy
```

## Pour tester votre projet :

### Construisez l'exécutable :
Ceci compile votre application et crée un fichier url-shortener à la racine du projet.
```bash
go build -o url-shortener
```
Désormais, toutes les commandes seront lancées avec ./url-shortener.

### Initialisation de la Base de Données

Avant de démarrer le serveur, créez le fichier de base de données SQLite et ses tables :

1.  **Exécutez les migrations :**
```bash
./url-shortener migrate
```
Un message de succès confirmera la création des tables. Un fichier url_shortener.db sera créé à la racine du projet.

### Lancer le Serveur et les Processus de Fond

C'est l'étape qui démarre le cœur de votre application. Elle démarre le serveur web, les workers qui enregistrent les clics, et le moniteur d'URLs.

Démarrez le service :
```bash
./url-shortener run-server
```
Laissez ce terminal ouvert et actif. Il affichera les logs du serveur HTTP, des workers de clics et du moniteur d'URLs.

### 4. Interagir avec le Service (Utilise un **Nouveau Terminal**)

Ouvre une **nouvelle fenêtre de terminal** pour exécuter les commandes CLI et tester les APIs pendant que le serveur est en cours d'exécution.

#### 4.1. Créer une URL courte (via la CLI)

Raccourcis une URL longue en utilisant la commande `create` :

```bash
./url-shortener create --url="https://www.example.com/ma-super-url-de-test-pour-le-tp-go-final"
```
Tu obtiendras un message similaire à :
```bash
URL courte créée avec succès:
Code: XYZ123
URL complète: http://localhost:8080/XYZ123
```

Note le Code (ex: XYZ123) et l'URL complète pour les étapes suivantes.

#### 4.2. Accéder à l'URL courte (via Navigateur)
1. Ouvre ton navigateur web et accède à l'URL complète que tu as obtenue (par exemple, http://localhost:8080/XYZ123).
2. Le navigateur devrait te rediriger instantanément vers l'URL longue originale. Dans le terminal où le serveur tourne (./url-shortener run-server), tu devrais voir des logs indiquant qu'un clic a été détecté et envoyé au worker asynchrone.

#### 4.3. Consulter les Statistiques (via la CLI)
go mod download
go mod tidy

# Compiler l'application
go build -o url-shortener
```

### 2. Configuration

Le fichier de configuration se trouve dans `configs/config.yaml` :

```yaml
# Configuration du serveur
server:
  port: 8080
  base_url: "http://localhost:8080"

# Configuration de la base de données
database:
  name: "url_shortener.db"

# Configuration des analytics
analytics:
  buffer_size: 1000
  worker_count: 5

# Configuration du moniteur
monitor:
  interval_minutes: 5
```

### 3. Initialiser la Base de Données

```bash
# Exécuter les migrations
./url-shortener migrate
```

### 4. Lancer le Serveur

```bash
# Démarrer le serveur API avec workers et moniteur
./url-shortener run-server
```

Le serveur démarre sur `http://localhost:8080`

## 📖 Utilisation

### Interface CLI

#### Créer une URL courte

```bash
# Création d'une URL courte
./url-shortener create --url="https://github.com/"

# Exemple de sortie :
# URL courte créée avec succès !
# Code court: kwDkxi
# URL complète: http://localhost:8080/kwDkxi
```

#### Consulter les statistiques

```bash
# Afficher les statistiques d'un lien
./url-shortener stats --code="kwDkxi"

# Exemple de sortie :
# Statistiques pour le code court: kwDkxi
# URL longue: https://github.com/
# Total de clics: 15
```

### API REST

#### 1. Vérifier l'état du service

```bash
curl --location 'http://localhost:8080/health'
```

**Réponse :**
```json
{
  "status": "ok"
}
```

#### 2. Créer une URL courte

```bash
curl --location 'http://localhost:8080/api/v1/links' \
--header 'Content-Type: application/json' \
--data '{"long_url":"https://www.google.com"}'
```

**Réponse :**
```json
{
  "short_code": "6Zc1qP",
  "short_url": "http://localhost:8080/6Zc1qP",
  "long_url": "https://www.google.com"
}
```

#### 3. Redirection vers l'URL originale

```bash
# Redirige automatiquement vers l'URL longue (HTTP 302)
curl --location 'http://localhost:8080/6Zc1qP'
```

Ou ouvrez simplement `http://localhost:8080/6Zc1qP` dans votre navigateur.

#### 4. Obtenir les statistiques d'un lien

```bash
curl --location 'http://localhost:8080/api/v1/links/6Zc1qP/stats'
```

**Réponse :**
```json
{
  "short_code": "6Zc1qP",
  "long_url": "https://www.google.com",
  "total_clicks": 42,
  "is_active": true
}
```

## 🏗️ Architecture du Projet

```
ProjetGo/
├── cmd/
│   ├── root.go              # Commande racine Cobra
│   ├── server/
│   │   └── server.go        # Lancement du serveur, workers et moniteur
│   └── cli/
│       ├── create.go        # Création d'URL via CLI
│       ├── stats.go         # Consultation des stats via CLI
│       └── migrate.go       # Migrations de la base de données
├── internal/
│   ├── api/
│   │   └── handlers.go      # Handlers HTTP (Gin)
│   ├── models/
│   │   ├── link.go          # Modèle GORM Link
│   │   └── click.go         # Modèle GORM Click
│   ├── services/
│   │   ├── link_service.go  # Logique métier des liens
│   │   └── click_service.go # Logique métier des clics
│   ├── workers/
│   │   └── click_worker.go  # Workers asynchrones pour analytics
│   ├── monitor/
│   │   └── url_monitor.go   # Monitoring périodique des URLs
│   ├── config/
│   │   └── config.go        # Configuration Viper
│   └── repository/
│       ├── link_repository.go  # Repository GORM pour Link
│       └── click_repository.go # Repository GORM pour Click
├── configs/
│   └── config.yaml          # Configuration du projet
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ Technologies Utilisées

- **[Go](https://go.dev/)** 1.24.3 - Langage de programmation
- **[Gin](https://gin-gonic.com/)** - Framework web HTTP
- **[GORM](https://gorm.io/)** - ORM pour Go
- **[SQLite](https://www.sqlite.org/)** - Base de données embarquée
- **[Cobra](https://cobra.dev/)** - CLI puissante
- **[Viper](https://github.com/spf13/viper)** - Gestion de configuration

## 🎯 Concepts Techniques

### Analytics Asynchrones
Les clics sont enregistrés en arrière-plan via :
- **Goroutines** : Workers dédiés à l'enregistrement
- **Channels bufferisés** : File d'attente des événements de clic
- **Non-bloquant** : La redirection ne dépend pas de l'enregistrement

### Monitoring d'URLs
- Vérification périodique de la disponibilité des URLs (HTTP 200/3xx)
- Notifications en cas de changement d'état (accessible ↔ inaccessible)
- Intervalle configurable via `config.yaml`

### Génération de Codes Courts
- Codes uniques de 6 caractères alphanumériques
- Gestion des collisions avec logique de retry
- Stockage persistant en base SQLite

## 📝 Exemples d'Utilisation Complets

### Scénario 1 : Création et utilisation via API

```bash
# 1. Créer une URL courte
curl --location 'http://localhost:8080/api/v1/links' \
--header 'Content-Type: application/json' \
--data '{"long_url":"https://www.google.com"}'

# Réponse : {"short_code":"abc123",...}

# 2. Utiliser l'URL courte (dans le navigateur ou via curl)
curl --location 'http://localhost:8080/abc123'

# 3. Consulter les stats
curl --location 'http://localhost:8080/api/v1/links/abc123/stats'
```

### Scénario 2 : Création et consultation via CLI

```bash
# 1. Créer une URL
./url-shortener create --url="https://github.com/"

# 2. Tester dans le navigateur
# Ouvrir : http://localhost:8080/kwDkxi

# 3. Consulter les statistiques
./url-shortener stats --code="kwDkxi"
```

## 🔧 Arrêt du Serveur

Pour arrêter proprement le serveur :
```
Ctrl + C
```

## 📚 Documentation Technique

### Endpoints API Détaillés

| Méthode | Endpoint | Description | Body/Params |
|---------|----------|-------------|-------------|
| GET | `/health` | Santé du service | - |
| POST | `/api/v1/links` | Créer URL courte | `{"long_url": "..."}` |
| GET | `/{shortCode}` | Redirection | - |
| GET | `/api/v1/links/{shortCode}/stats` | Statistiques | - |

### Commandes CLI Détaillées

| Commande | Description | Options |
|----------|-------------|---------|
| `run-server` | Lance le serveur | - |
| `create` | Crée une URL courte | `--url` (requis) |
| `stats` | Affiche les stats | `--code` (requis) |
| `migrate` | Migrations DB | - |

## 👨‍💻 Développement

### Structure des Commits
- Messages clairs et descriptifs
- Organisation logique des changements
- Respect des conventions Git

### Qualité du Code
- Code commenté et documenté
- Respect des conventions Go
- Gestion d'erreurs pertinente
- Architecture propre (Repository, Service patterns)

## 📄 Licence

Projet développé dans le cadre d'un TP Go Final.

---

**Auteur** : [axellelanca](https://github.com/axellelanca)  
**Date** : 2025
