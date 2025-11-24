package repository

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

// LinkRepository est une interface qui définit les méthodes d'accès aux données
// pour les opérations CRUD sur les liens.
type LinkRepository interface {
	CreateLink(link *models.Link) error
	GetLinkByShortCode(shortCode string) (*models.Link, error)
	GetAllLinks() ([]models.Link, error)
	CountClicksByLinkID(linkID uint) (int, error)
}

// GormLinkRepository est l'implémentation de LinkRepository utilisant GORM.
type GormLinkRepository struct {
	db *gorm.DB // Référence à l'instance de la base de données GORM
}

// NewLinkRepository crée et retourne une nouvelle instance de GormLinkRepository.
// Cette fonction retourne *GormLinkRepository, qui implémente l'interface LinkRepository.
func NewLinkRepository(db *gorm.DB) *GormLinkRepository {
	return &GormLinkRepository{db: db}
}

// CreateLink insère un nouveau lien dans la base de données.
func (r *GormLinkRepository) CreateLink(link *models.Link) error {
	// Utilise la méthode Create de GORM pour insérer le lien dans la table 'links'
	result := r.db.Create(link)
	if result.Error != nil {
		return fmt.Errorf("failed to create link: %w", result.Error)
	}
	return nil
}

// GetLinkByShortCode récupère un lien de la base de données en utilisant son shortCode.
// Il renvoie gorm.ErrRecordNotFound si aucun lien n'est trouvé avec ce shortCode.
func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	// Utilise la méthode First de GORM pour trouver le premier lien avec le ShortCode donné
	result := r.db.Where("short_code = ?", shortCode).First(&link)
	if result.Error != nil {
		return nil, result.Error // Retourne l'erreur (peut être gorm.ErrRecordNotFound)
	}
	return &link, nil
}

// GetAllLinks récupère tous les liens de la base de données.
// Cette méthode est utilisée par le moniteur d'URLs.
func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) {
	var links []models.Link
	// Utilise la méthode Find de GORM pour récupérer tous les enregistrements de la table 'links'
	result := r.db.Find(&links)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve all links: %w", result.Error)
	}
	return links, nil
}

// CountClicksByLinkID compte le nombre total de clics pour un ID de lien donné.
func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64 // GORM retourne un int64 pour les comptes
	// Utilise la méthode Count de GORM pour compter les clics où LinkID correspond à l'ID donné
	result := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count clicks for link ID %d: %w", linkID, result.Error)
	}
	return int(count), nil
}
