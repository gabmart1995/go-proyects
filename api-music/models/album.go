package models

import "time"

type Album struct {
	Id          int `json:"id" gorm:"primaryKey,autoIncrement"`
	UserId      int
	Name        string    `json:"name" validate:"required,min=2,max=30"`
	Description string    `json:"description" validate:"required,min=2,max=30"`
	PublishDate time.Time `json:"publish_date" validate:"required,min=2,max=30"`
	Image       string    `json:"image"`

	// campos que representan las relaciones
	Artist Artist `gorm:"foreignKey:AlbumId" validate:"omitempty"`
	Songs  []Song `gorm:"foreignKey:AlbumId" validate:"omitempty"` // relacion one to many
}
