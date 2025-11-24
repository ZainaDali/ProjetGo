package models

// TODO : Créer la struct Link
// Link représente un lien raccourci dans la base de données.
// Les tags `gorm:"..."` définissent comment GORM doit mapper cette structure à une table SQL.
// ID qui est une primaryKey
// Shortcode : doit être unique, indexé pour des recherches rapide (voir doc), taille max 10 caractères
// LongURL : doit pas être null
// CreateAt : Horodatage de la créatino du lien


// Link représente un lien raccourci dans la base de données.
type Link struct {
	ID        uint      `gorm:"primaryKey"`                    // Clé primaire
	ShortCode string    `gorm:"size:10;uniqueIndex;not null"`  // Code court unique, indexé et max 10 caractères
	LongURL   string    `gorm:"not null"`                      // Lien original obligatoire
	CreatedAt time.Time                                       // Horodatage de la création du lien
}
