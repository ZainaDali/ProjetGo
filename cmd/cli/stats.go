package cli

import (
	"fmt"
	"log"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO : variable shortCodeFlag qui stockera la valeur du flag --code
var shortCodeFlag string

// StatsCmd représente la commande 'stats'
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,

	Run: func(cmd *cobra.Command, args []string) {

		// TODO : Valider que le flag --code a été fourni.
		if shortCodeFlag == "" {
			log.Println("ERREUR : Le flag --code est obligatoire.")
			os.Exit(1)
		}

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Println("ERREUR : Impossible de charger la configuration globale.")
			os.Exit(1)
		}

		// TODO : Initialiser la connexion à la BDD.
		db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
		if err != nil {
			log.Fatalf("ERREUR : Impossible d'ouvrir la base SQLite: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL : Impossible d'obtenir la base SQL sous-jacente : %v", err)
		}

		// TODO : fermer la connexion avec defer
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// TODO 5: Appeler GetLinkStats
		link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("Aucun lien trouvé pour le code : %s\n", shortCodeFlag)
				os.Exit(1)
			}
			log.Printf("ERREUR : Impossible de récupérer les statistiques : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Statistiques pour le code court: %s\n", link.ShortCode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n", totalClicks)
	},
}

// init() définit les flags et ajoute la commande
func init() {

	// TODO : Définir le flag --code pour la commande stats.
	StatsCmd.Flags().StringVar(&shortCodeFlag, "code", "", "Code court de l'URL")

	// TODO : Marquer le flag comme requis
	StatsCmd.MarkFlagRequired("code")

	// TODO : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(StatsCmd)
}

