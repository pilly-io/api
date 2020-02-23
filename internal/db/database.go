package db

import "time"

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
	Conditions     QueryConditions
	Interval       *QueryInterval
	Limit          int
	OrderBy        string
	Page           int
	ExcludeDeleted bool
}

// Database interface
type Database interface {
	Migrate()
	Clusters() *ClustersTable
	Nodes() Table
	Metrics() *MetricsTable
	Owners() *OwnersTable
	Namespaces() Table
	Flush()
}
