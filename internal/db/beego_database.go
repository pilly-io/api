package db

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pilly-io/api/internal/models"
)

type BeegoDatabase struct {
	client *orm.Ormer

	clusters   *ClustersTable
	nodes      Table
	metrics    *MetricsTable
	owners     *OwnersTable
	namespaces *NamespacesTable
}

const dbAlias = "default"

func init() {
	orm.RegisterModel(new(models.Cluster), new(models.Node), new(models.Namespace), new(models.Metric), new(models.Owner))
}

func NewBeegoDatabase(uri string) *BeegoDatabase {
	orm.RegisterDriver("postgres", orm.DRPostgres)

	orm.RegisterDataBase(dbAlias, "postgres", uri)

	client := orm.NewOrm()
	clusters := NewClusterTable(client, models.Cluster{})
	nodes := NewBeegoTable(client, models.Node{})
	metrics := NewMetricsTable(client, models.Metric{})
	namespaces := NewNamespacesTable(client, models.Namespace{})
	owners := NewOwnersTable(client, models.Owner{})
	return &BeegoDatabase{
		&client,
		clusters,
		nodes,
		metrics,
		owners,
		namespaces,
	}
}

// Migrate sync the schemas of the DB
func (db *BeegoDatabase) Migrate() {
	orm.RunSyncdb(dbAlias, true, false)
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
func (db *BeegoDatabase) Namespaces() *NamespacesTable {
	return db.namespaces
}

// Flush delete all records from all tables
func (db *BeegoDatabase) Flush() {
	db.Clusters().Delete(Query{}, false)
	db.Namespaces().Delete(Query{}, false)
	db.Owners().Delete(Query{}, false)
	db.Metrics().Delete(Query{}, false)
	db.Nodes().Delete(Query{}, false)
}
