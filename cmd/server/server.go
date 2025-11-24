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
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// RunServerCmd représente la commande 'run-server' de Cobra.
var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
démarre les workers asynchrones pour les clics et le moniteur d'URLs,
puis lance le serveur HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {

		// TODO : créer une variable qui stock la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("FATAL: Impossible de charger la configuration globale.")
		}

		// TODO : Initialiser la connexion à la bBDD
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Impossible d'ouvrir la base SQLite: %v", err)
		}

		// Initialiser les repositories
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)

	log.Println("Repositories initialisés.")

	// Initialiser les services métiers
	linkService := services.NewLinkService(linkRepo)

	log.Println("Services métiers initialisés.")

	// TODO : Initialiser le channel ClickEventsChannel + workers
	clickChan := make(chan models.ClickEvent, cfg.Analytics.BufferSize)

	workers.StartClickWorkers(cfg.Analytics.WorkerCount, clickChan, clickRepo)

	log.Printf(
		"Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
		cfg.Analytics.BufferSize, cfg.Analytics.WorkerCount,
	)

	// TODO : Initialiser et lancer le moniteur d'URLs
	monitorInterval := time.Duration(cfg.Monitor.IntervalMinutes) * time.Minute
	urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)

	go urlMonitor.Start()

	log.Printf("Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)

	// TODO : Configurer le routeur Gin et les handlers API.
	router := gin.Default()
	api.SetupRoutes(router, linkService)

	log.Println("Routes API configurées.")		// Créer le serveur HTTP Gin
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// TODO : Démarrer le serveur Gin dans une goroutine
		go func() {
			log.Printf("Serveur lancé sur %s...", serverAddr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Erreur serveur HTTP: %v", err)
			}
		}()

		// TODO : Gérer l'arrêt propre du serveur
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit
		log.Println("Signal d'arrêt reçu. Arrêt du serveur...")

		// Attente pour laisser les workers finir
		log.Println("Arrêt en cours... Donnez un peu de temps aux workers pour finir.")
		time.Sleep(5 * time.Second)

		log.Println("Serveur arrêté proprement.")
	},
}

func init() {
	// TODO : ajouter la commande
	cmd2.RootCmd.AddCommand(RunServerCmd)
}
