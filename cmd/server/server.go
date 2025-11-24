package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/api"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/monitor"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/axellelanca/urlshortener/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	sqlite "github.com/glebarez/sqlite" // DRIVER 100% GO
	"gorm.io/gorm"
)

var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API et les processus de fond.",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("FATAL : Impossible de charger la configuration.")
		}

		// Connexion SQLite SANS modernc
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL : Impossible d'ouvrir SQLite : %v", err)
		}

		// Repos
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)

		log.Println("Repositories initialisés.")

		// Services
		linkService := services.NewLinkService(linkRepo)

		log.Println("Services métiers initialisés.")

		// Channel + workers
		clickChan := make(chan models.ClickEvent, cfg.Analytics.BufferSize)
		workers.StartClickWorkers(cfg.Analytics.WorkerCount, clickChan, clickRepo)

		log.Printf("Channel clic prêt : buffer=%d workers=%d",
			cfg.Analytics.BufferSize, cfg.Analytics.WorkerCount)

		// Moniteur d'URL
		monitorInterval := time.Duration(cfg.Monitor.IntervalMinutes) * time.Minute
		urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)
		go urlMonitor.Start()

		log.Printf("Monitor démarré (%v)", monitorInterval)

		// Routes
		router := gin.Default()
		api.SetupRoutes(router, linkService, clickChan)

		log.Println("Routes API configurées.")

		// Serveur HTTP
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

        // Serveur asynchrone
		go func() {
			log.Printf("Serveur lancé sur %s ...", serverAddr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Erreur serveur : %v", err)
			}
		}()

		// Shutdown propre
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("Arrêt demandé. Fermeture...")

		time.Sleep(5 * time.Second)

		log.Println("Serveur arrêté.")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(RunServerCmd)
}
