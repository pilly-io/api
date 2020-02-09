package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pilly-io/api/internal/models"
)

type PaginationInfo struct {
	CurrentPage int
	MaxPage     int
	TotalCount  int
	Limit       int
}

type QueryConditions = map[string]interface{}

type Query struct {
	Conditions QueryConditions
	Limit      int
	OrderBy    string
	Desc       bool
	Page       int
}

// GormDatabase is wrapper for the orm
type GormDatabase struct {
	*gorm.DB
	clusters *ClustersTable
	nodes    Table
}

// Database interface
type Database interface {
	Migrate()
	Insert(value interface{}) error
	Clusters() *ClustersTable
	Nodes() Table
}

// New creates an new DB object
func New(driver string, DBURI string) (*GormDatabase, error) {
	db, err := gorm.Open(driver, DBURI)
	db.LogMode(true)
	return &GormDatabase{
		db,
		NewClusterTable(db, models.Cluster{}),
		NewTable(db, models.Node{}),
	}, err
}

// Migrate sync the schemas of the DB
func (db *GormDatabase) Migrate() {
	db.AutoMigrate(&models.Cluster{})
}

// Insert creates a new record in the right table
func (db *GormDatabase) Insert(value interface{}) error {
	return db.Create(value).Error
}

// Clusters returns the clusters Table object
func (db *GormDatabase) Clusters() *ClustersTable {
	return db.clusters
}

// Nodes returns the nodes Table object
func (db *GormDatabase) Nodes() Table {
	return db.nodes
}
