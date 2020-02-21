package db

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

type QueryInterval struct {
	Start time.Time
	End   time.Time
}

type Query struct {
	Conditions QueryConditions
	Interval   *QueryInterval
	Limit      int
	OrderBy    string
	Desc       bool
	Page       int
}

// GormDatabase is wrapper for the orm
type GormDatabase struct {
	*gorm.DB
	clusters   *ClustersTable
	Node       Table
	metrics    *MetricsTable
	owners     *OwnersTable
	namespaces Table
}

// Database interface
type Database interface {
	Migrate()
	Insert(value interface{}) error
	Clusters() *ClustersTable
	Nodes() Table
	Metrics() *MetricsTable
	Owners() *OwnersTable
	Namespaces() *Table
	Flush()
}

// getAllModels returns all the models used in this DB, used to migrate DB and truncate
// (can't use a const because go does not support array as constant)
func getAllModels() []interface{} {
	return []interface{}{&models.Cluster{}, &models.Node{}, &models.Owner{}, models.Metric{}, models.Namespace{}}
}

// New creates an new DB object
func New(driver string, DBURI string) (*GormDatabase, error) {
	db, err := gorm.Open(driver, DBURI)
	db.LogMode(true)
	return &GormDatabase{
		db,
		NewClusterTable(db, models.Cluster{}),
		NewTable(db, models.Node{}),
		NewMetricsTable(db, models.Metric{}),
		NewOwnersTable(db, models.Owner{}),
		NewTable(db, models.Namespace{}),
	}, err
}

// Migrate sync the schemas of the DB
func (db *GormDatabase) Migrate() {
	//db.AutoMigrate(&models.Cluster{}, &models.Node{}, &models.Namespace{}, &models.Owner{}, models.Metric{})
	db.AutoMigrate(getAllModels()...)
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
	return db.Node
}

// Metrics returns the metrics Table object
func (db *GormDatabase) Metrics() *MetricsTable {
	return db.metrics
}

// Owners returns the owners Table object
func (db *GormDatabase) Owners() *OwnersTable {
	return db.owners
}

// Namespaces returns the namespaces Table object
func (db *GormDatabase) Namespaces() *Table {
	return db.namespaces
}

// Flush remove all records from all tables
func (db *GormDatabase) Flush() {
	for _, model := range getAllModels() {
		db.Unscoped().Delete(model)
	}
}
