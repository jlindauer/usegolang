package models

import "github.com/jinzhu/gorm"

type Gallery struct {
  gorm.Model
  UserID uint  `gorm:"not_null;index"`
  Title string `gorm:"not_null"`
}

type galleryService struct {
  GalleryDB
}

type galleryValidator struct {
  GalleryDB
}

var _ GalleryDB = &galleryGorm{}

type GalleryService interface {
  GalleryDB
}

type GalleryDB interface {
  Create(gallery *Gallery) error
}

type galleryGorm struct {
  db *gorm.DB
}

func NewGalleryService(db *gorm.DB) GalleryService {
  return &galleryService{
    GalleryDB: &galleryValidator{
      GalleryDB: &galleryGorm{
        db: db,
      },
    },
  }
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
  // TODO: Implement this later
  return nil
}
