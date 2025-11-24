package cmd

import (
	"log"

	"github.com/axellelanca/urlshortener/internal/config"
	"github.com/spf13/cobra"
)

// Cfg est la configuration globale chargée via Viper.
var Cfg *config.Config

// RootCmd représente la commande racine.
var RootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
	Long: `'url-shortener' est une application complète pour gérer des URLs courtes.
Elle inclut un serveur API pour le raccourcissement et la redirection,
ainsi qu'une interface en ligne de commande pour l'administration.

Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.`,
}

// Execute exécute la commande racine.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// init s'exécute automatiquement avant main().
func init() {
	// Charger la configuration au lancement de n'importe quelle commande
	cobra.OnInitialize(initConfig)
}

// initConfig charge la configuration de l'application via Viper.
func initConfig() {
	var err error
	Cfg, err = config.LoadConfig()

	if err != nil {
		log.Printf("Attention: Problème lors du chargement de la configuration: %v. Utilisation des valeurs par défaut.", err)
	}
}

