package models

import (
	"database/sql"
	"time"
)

type Dataset struct {
	ModelByInt
	Name          string
	Namespace     string
	Description   string
	Preprocessing string
	Instances     int
	Format        string
	Labels        []DatasetLabel
	PublishAt     time.Time
	ModifyAT      time.Time
	Resources     []DatasetResource
	Lisense       string
	UserID        uint
	User          User
	Private       sql.NullBool
	// Reference     []string
	Creator string
}

type DatasetLabel struct {
	ModelByInt
	Name    string
	Summary string
	Color   int
}

type DatasetResource struct {
	ModelByInt
	Storage DatasetStorage
	Type    string
	Path    string
}

type DatasetStorage struct {
	ModelByInt
	AccessID  string
	AccessKey string
	Region    string
	Bucket    string
	Endpoint  string
	// ACL: awss3.BucketCannedACLPublicRead,

}

func GetDataset(name string) *Dataset {
	var ds Dataset
	if db.Find(&ds, "name = ? and namespace IS NULL", name).RecordNotFound() {
		return nil
	}
	return &ds
}

func GetDatasetByNamespace(name, ns string) *Dataset {
	var ds Dataset
	if db.Find(&ds, "name = ? and namespace = ?", name, ns).RecordNotFound() {
		return nil
	}
	return &ds
}
