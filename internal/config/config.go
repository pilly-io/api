package config

import "github.com/jinzhu/gorm"

var Settings Config

// note: maybe we should do like https://github.com/gin-gonic/gin/issues/932#issuecomment-487297482 . IDC
type Config struct {
	// the shared DB ORM object
	DB *gorm.DB
	// the error thrown be GORM when using DB ORM object
	DBErr error
}
