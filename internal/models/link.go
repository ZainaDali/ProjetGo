package models

import "time"

// Link représente un lien raccourci dans la base de données.
// Les tags `gorm:"..."` définissent comment GORM doit mapper cette structure à une table SQL.
type Link struct {
	ID        uint      `gorm:"primaryKey"`                    // ID est la clé primaire
	ShortCode string    `gorm:"unique;index;size:10;not null"` // Shortcode : doit être unique, indexé pour des recherches rapides, taille max 10 caractères
	LongURL   string    `gorm:"not null"`                      // LongURL : ne doit pas être null
	CreatedAt time.Time // Horodatage de la création du lien
}
