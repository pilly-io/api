package db

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pilly-io/api/internal/models"
)

type BeegoDatabase struct {
	client *orm.Ormer

	clusters   *ClustersTable
	Node       Table
	metrics    *MetricsTable
	owners     *OwnersTable
	namespaces Table
	tables []Table
}

const dbAlias = "default"

func NewBeegoDatabase(uri string) *BeegoDatabase {
	orm.RegisterDriver("postgres", orm.DRPostgres)

	orm.RegisterDataBase(dbAlias, "postgres", uri)
	orm.RegisterModel(new(models.Cluster), new(models.Node), new(models.Namespace), new(models.Metric), new(models.Owner))

	client := orm.NewOrm()
	var clusters, nodes, metrics, owners, namespaces Table
	return &BeegoDatabase{
		client: &client,
		tables: []Table{clusters, nodes, metrics, owners, namespaces}
	}
}

// Migrate sync the schemas of the DB
func (db *BeegoDatabase) Migrate() {
	//db.AutoMigrate(&models.Cluster{}, &models.Node{}, &models.Namespace{}, &models.Owner{}, models.Metric{})
	orm.RunSyncdb(dbAlias, false, true)
}


// Clusters returns the clusters Table object
func (db *BeegoDatabase) Clusters() *ClustersTable {
	return db.clusters
}

// Nodes returns the nodes Table object
func (db *BeegoDatabase) Nodes() Table {
	return db.nodes
}

// Metrics returns the metrics Table object
func (db *BeegoDatabase) Metrics() *MetricsTable {
	return db.metrics
}

// Owners returns the owners Table object
func (db *BeegoDatabase) Owners() *OwnersTable {
	return db.owners
}

// Namespaces returns the namespaces Table object
func (db *BeegoDatabase) Namespaces() Table {
	return db.namespaces
}

// Flush delete all records from all tables
func (db *BeegoDatabase) Flush() {
	for _, table := range db.tables {
		table.Delete(Query{})
	}
}
