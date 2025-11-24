package repository

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

type ClickRepository interface {
	CreateClick(click *models.Click) error
	CountClicksByLinkID(linkID uint) (int, error)
}

type GormClickRepository struct {
	db *gorm.DB
}

func NewClickRepository(db *gorm.DB) *GormClickRepository {
	return &GormClickRepository{db: db}
}

func (r *GormClickRepository) CreateClick(click *models.Click) error {
	result := r.db.Create(click)
	if result.Error != nil {
		return fmt.Errorf("failed to create click: %w", result.Error)
	}
	return nil
}

func (r *GormClickRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	result := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count clicks for link ID %d: %w", linkID, result.Error)
	}
	return int(count), nil
}
