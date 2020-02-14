package models

import (
	"time"
)

//Model : a copy of gorm.Model with json annotations
type Model struct {
	ID        uint       `gorm:"primary_key;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type Cluster struct {
	Model
	Name       string `gorm:"unique;not null" json:"name"`
	Provider   string `json:"provider"`
	Region     string `json:"region"`
	APIToken   string `json:"api_token"`
	NodesCount int    `json:"nodes_count"`
}

type Node struct {
	Model
	InstanceType      string            `json:"instance_type"`
	Region            string            `json:"region"`
	Zone              string            `json:"zone"`
	Hostname          string            `json:"hostname"`
	UID               string            `json:"uid"`
	KubernetesVersion string            `json:"kubernetes_version"`
	OS                string            `json:"os"`
	ClusterID         int               `json:"cluster_id"`
	Labels            map[string]string `json:"labels"`
}

type Namespace struct {
	Model
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	ClusterID int               `json:"cluster_id"`
}

type Metric struct {
	Model
	Name      string `json:"metric_name"`
	Value     string `json:"metric_value"`
	OwnerUID  string `json:"owner_uid"`
	ClusterID int    `json:"cluster_id"`
}

type Owner struct {
	Model
	UID       string            `json:"uid"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	ClusterID int               `json:"cluster_id"`
}
