package repository

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

type LinkRepository interface {
	CreateLink(link *models.Link) error
	GetLinkByShortCode(shortCode string) (*models.Link, error)
	GetAllLinks() ([]models.Link, error)
	CountClicksByLinkID(linkID uint) (int, error)
}

type GormLinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *GormLinkRepository {
	return &GormLinkRepository{db: db}
}

func (r *GormLinkRepository) CreateLink(link *models.Link) error {
	result := r.db.Create(link)
	if result.Error != nil {
		return fmt.Errorf("failed to create link: %w", result.Error)
	}
	return nil
}

func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	result := r.db.Where("short_code = ?", shortCode).First(&link)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) {
	var links []models.Link
	result := r.db.Find(&links)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve all links: %w", result.Error)
	}
	return links, nil
}

func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	result := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count clicks for link ID %d: %w", linkID, result.Error)
	}
	return int(count), nil
}
