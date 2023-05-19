package migrations

import (
	"api-music/models"

	"gorm.io/gorm"
)

func Execute(db *gorm.DB) {

	migrator := db.Migrator()

	// -- tables -- //
	// users
	if !migrator.HasTable(&models.User{}) {
		migrator.CreateTable(&models.User{})
	}

	// album
	if !migrator.HasTable(&models.Album{}) {
		migrator.CreateTable(&models.Album{})
	}

	// songs
	if !migrator.HasTable(&models.Song{}) {
		migrator.CreateTable(&models.Song{})
	}

	// artist
	if !migrator.HasTable(&models.Artist{}) {
		migrator.CreateTable(&models.Artist{})
	}

	// -- constraints -- //
	// foreign keys album-song
	if !migrator.HasConstraint(&models.Album{}, "User") &&
		!migrator.HasConstraint(&models.Album{}, "fk_album_user") {

		migrator.CreateConstraint(&models.Album{}, "User")
		migrator.CreateConstraint(&models.Album{}, "fk_album_user")
	}

	// foreign keys album-artist
	if !migrator.HasConstraint(&models.Album{}, "Artist") &&
		!migrator.HasConstraint(&models.Album{}, "fk_album_artist") {

		migrator.CreateConstraint(&models.Album{}, "Artist")
		migrator.CreateConstraint(&models.Album{}, "fk_album_artist")
	}

	if !migrator.HasConstraint(&models.Album{}, "Songs") &&
		!migrator.HasConstraint(&models.Album{}, "fk_album_songs") {

		migrator.CreateConstraint(&models.Album{}, "Songs")
		migrator.CreateConstraint(&models.Album{}, "fk_album_songs")
	}
}
