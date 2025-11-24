package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// Définition du jeu de caractères pour la génération des codes courts.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// LinkService est une structure qui fournit des méthodes pour la logique métier des liens.
// Elle détient linkRepo qui est une référence vers une interface LinkRepository.
type LinkService struct {
	linkRepo repository.LinkRepository // IMPORTANT : Le champ est du type de l'interface (non-pointeur)
}

// NewLinkService crée et retourne une nouvelle instance de LinkService.
func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

// GenerateShortCode est une méthode rattachée à LinkService.
// Elle génère un code court aléatoire d'une longueur spécifiée.
// Il utilise le package 'crypto/rand' pour éviter la prévisibilité.
func (s *LinkService) GenerateShortCode(length int) (string, error) {
	// Crée un slice de bytes pour stocker le code généré
	code := make([]byte, length)

	// Génère chaque caractère du code de manière aléatoire
	for i := range code {
		// Génère un nombre aléatoire entre 0 et len(charset)-1
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random short code: %w", err)
		}
		// Assigne le caractère correspondant au code
		code[i] = charset[randomIndex.Int64()]
	}

	return string(code), nil
}


// CreateLink crée un nouveau lien raccourci.
// Il génère un code court unique, puis persiste le lien dans la base de données.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	// Variable pour stocker le shortcode créé
	var shortCode string

	// Définir un nombre maximum de tentatives pour trouver un code unique
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		// Génère un code de 6 caractères
		code, err := s.GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		// Vérifie si le code généré existe déjà en base de données
		_, err = s.linkRepo.GetLinkByShortCode(code)
		if err != nil {
			// Si l'erreur est 'record not found' de GORM, cela signifie que le code est unique.
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code // Le code est unique, on peut l'utiliser
				break            // Sort de la boucle de retry
			}
			// Si c'est une autre erreur de base de données, retourne l'erreur.
			return nil, fmt.Errorf("database error checking short code uniqueness: %w", err)
		}

		// Si aucune erreur (le code a été trouvé), cela signifie une collision.
		log.Printf("Short code '%s' already exists, retrying generation (%d/%d)...", code, i+1, maxRetries)
		// La boucle continuera pour générer un nouveau code.
	}

	// Si après toutes les tentatives, aucun code unique n'a été trouvé
	if shortCode == "" {
		return nil, errors.New("failed to generate a unique short code after multiple attempts")
	}

	// Crée une nouvelle instance du modèle Link
	link := &models.Link{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}

	// Persiste le nouveau lien dans la base de données via le repository
	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("failed to create link in database: %w", err)
	}

	// Retourne le lien créé
	return link, nil
}

// GetLinkByShortCode récupère un lien via son code court.
// Il délègue l'opération de recherche au repository.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	// Récupère le lien par son code court en utilisant le repository
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, err
	}
	return link, nil
}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
// Il interagit avec le LinkRepository pour obtenir le lien, puis avec le ClickRepository
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// Récupère le lien par son shortCode
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, err
	}

	// Compte le nombre de clics pour ce LinkID
	totalClicks, err := s.linkRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count clicks: %w", err)
	}

	// Retourne le lien, le nombre de clics et aucune erreur
	return link, totalClicks, nil
}

