package models

import (
  "fmt"
  "io"
  "os"
  "path/filepath"
)

type Image struct {
  GalleryID uint
  Filename  string
}

// Path is used to build the absolute path used to reference
// this image via a web request.
func (i *Image) Path() string {
  return "/" + i.RelativePath()
}

// RelativePath is used to build the path to this image on our local disk,
// relative to where our Go application is ran from.
func (i *Image) RelativePath() string {
  // Convert the gallery ID to a string
  galleryID := fmt.Sprintf("%v", i.GalleryID)
  return filepath.ToSlash(filepath.Join("images", "galleries", galleryID, i.Filename))
}

type ImageService interface {
  Create(galleryID uint, r io.Reader, filename string) error
  ByGalleryID(galleryID uint) ([]Image, error)
  Delete (i *Image) error
}

func NewImageService() ImageService {
  return &imageService{}
}

type imageService struct{}

func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {
  path, err := is.mkImageDir(galleryID)
  if err != nil {
    return err
  }
  // Create a destination file
  dst, err := os.Create(filepath.Join(path, filename))
  if err != nil {
    return err
  }
  defer dst.Close()
  // Copy reader data to the destination file
  _, err = io.Copy(dst, r)
  if err != nil {
    return err
  }
  return nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
  path := is.imageDir(galleryID)
  strings, err := filepath.Glob(filepath.Join(path, "*"))
  if err != nil {
    return nil, err
  }
  // Setup the Image slice being returned
  ret := make([]Image, len(strings))
  for i, imgStr := range strings {
    ret[i] = Image{
      Filename:  filepath.Base(imgStr),
      GalleryID: galleryID,
    }
  }
  return ret, nil
}

func (is *imageService) Delete(i *Image) error {
  return os.Remove(i.RelativePath())
}

func (is *imageService) imageDir(galleryID uint) string {
  return filepath.Join("images", "galleries", fmt.Sprintf("%v", galleryID))
}

func (is *imageService) mkImageDir(galleryID uint) (string, error) {
  galleryPath := is.imageDir(galleryID)
  err := os.MkdirAll(galleryPath, 0755)
  if err != nil {
    return "", err
  }
  return galleryPath, nil
}
