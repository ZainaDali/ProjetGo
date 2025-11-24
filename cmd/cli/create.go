package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"

	sqlite "github.com/glebarez/sqlite" // DRIVER SQLITE 100% Go (pas de CGO)
	"gorm.io/gorm"
)

// stocke la valeur du flag --url
var longURLFlag string

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,

	Run: func(cmd *cobra.Command, args []string) {

		// Vérifier que le flag --url est fourni
		if longURLFlag == "" {
			log.Println("ERREUR : Le flag --url est obligatoire.")
			os.Exit(1)
		}

		// Valider le format de l'URL
		if _, err := url.ParseRequestURI(longURLFlag); err != nil {
			log.Printf("ERREUR : Format d'URL invalide : %v\n", err)
			os.Exit(1)
		}

		// Charger la configuration globale
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Println("ERREUR : Impossible de charger la configuration globale.")
			os.Exit(1)
		}

		// Connexion SQLite via glebarez/sqlite (sans CGO)
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("ERREUR : Impossible d'ouvrir la base SQLite : %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("ERREUR FATALE : Impossible d'obtenir la base SQL sous-jacente : %v", err)
		}
		defer sqlDB.Close()

		// Initialiser repository + service
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// Créer le lien court
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			log.Printf("ERREUR : Impossible de créer le lien court : %v\n", err)
			os.Exit(1)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)

		fmt.Println("URL courte créée avec succès :")
		fmt.Printf("Code : %s\n", link.ShortCode)
		fmt.Printf("URL complète : %s\n", fullShortURL)
	},
}

func init() {

	// Définir le flag --url
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")

	// Rendre le flag obligatoire
	CreateCmd.MarkFlagRequired("url")

	// Ajouter la commande au root
	cmd2.RootCmd.AddCommand(CreateCmd)
}
