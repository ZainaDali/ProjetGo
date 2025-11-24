package cli

import (
	"fmt"
	"log"
	"net/url" // Pour valider le format de l'URL
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
	 
	 

)

// TODO : Faire une variable longURLFlag qui stockera la valeur du flag --url
var longURLFlag string

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,

	Run: func(cmd *cobra.Command, args []string) {

		// TODO 1: Valider que le flag --url a été fourni.
		if longURLFlag == "" {
			log.Println("ERREUR : Le flag --url est obligatoire.")
			os.Exit(1)
		}

		// TODO Validation basique du format de l'URL avec url.ParseRequestURI
		if _, err := url.ParseRequestURI(longURLFlag); err != nil {
			log.Printf("ERREUR : Format d'URL invalide: %v\n", err)
			os.Exit(1)
		}

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Println("ERREUR : Configuration introuvable.")
			os.Exit(1)
		}

		// TODO : Initialiser la connexion à la base de données SQLite.
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("ERREUR : Impossible d'ouvrir la base SQLite: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// TODO : Appeler CreateLink pour créer le lien court.
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			log.Printf("ERREUR : Impossible de créer le lien court: %v\n", err)
			os.Exit(1)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)

		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() définit les flags
func init() {

	// TODO : Définir le flag --url pour la commande create.
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")

	// TODO : Marquer le flag comme requis
	CreateCmd.MarkFlagRequired("url")

	// TODO : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(CreateCmd)
}

