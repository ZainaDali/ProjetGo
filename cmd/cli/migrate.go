package cli

import (
	"fmt"
	"log"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/spf13/cobra"

	sqlite "github.com/glebarez/sqlite" // DRIVER 100% GO (PAS modernc)
	"gorm.io/gorm"
)

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations pour créer ou mettre à jour les tables.",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("ERREUR : Impossible de charger la configuration globale.")
		}

		// Ouverture DB SANS CGO - SANS modernc
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("ERREUR : Impossible d'ouvrir la base : %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("ERREUR : Impossible d'obtenir la connexion SQL : %v", err)
		}
		defer sqlDB.Close()

		// Migrations GORM
		if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
			log.Fatalf("ERREUR : Migrations échouées : %v", err)
		}

		fmt.Println("Migrations exécutées avec succès.")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(MigrateCmd)
}


