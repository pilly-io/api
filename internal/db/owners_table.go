package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pilly-io/api/internal/models"
)

// OwnersTable implements Cluster specific operations
type OwnersTable struct {
	*GormTable
}

// NewOwnersTable : returns a new Table object
func NewOwnersTable(client *gorm.DB, model models.Owner) *OwnersTable {
	table := GormTable{client, model}
	return &OwnersTable{&table}
}
