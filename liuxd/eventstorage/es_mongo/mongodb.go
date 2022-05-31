/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package es_mongo

// mongodb package is an implementation of StateStore interface to perform operations on store

import (
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/common"
)

const (
	eventCollectionName     = "eventCollectionName"
	snapshotCollectionName  = "snapshotCollectionName"
	aggregateCollectionName = "aggregateCollectionName"
	id                      = "_id"
	value                   = "value"
	etag                    = "_etag"

	defaultEventCollectionName     = "dapr_event"
	defaultSnapshotCollectionName  = "dapr_snapshot"
	defaultAggregateCollectionName = "dapr_aggregate"
)

// MongoDB is a state store implementation for MongoDB.
type MongoDB struct {
	*common.MongoDB
	storageMetadata *storageMetadata
}

type storageMetadata struct {
	*common.MongoDBMetadata
	aggregateCollectionName string
	eventCollectionName     string
	snapshotCollectionName  string
}

// NewMongoDB returns a new MongoDB state store.
func NewMongoDB(logger logger.Logger) *MongoDB {
	mdb := common.NewMongoDB(logger)
	s := &MongoDB{
		MongoDB: mdb,
	}
	return s
}

// Init establishes connection to the store based on the metadata.
func (m *MongoDB) Init(metadata common.Metadata) error {
	if err := m.MongoDB.Init(metadata); err != nil {
		return err
	}
	storageMetadata, err := m.getStorageMetadata(metadata)
	if err != nil {
		return err
	}
	m.storageMetadata = storageMetadata
	return nil
}

func (m *MongoDB) getStorageMetadata(metadata common.Metadata) (*storageMetadata, error) {
	meta := storageMetadata{
		MongoDBMetadata:         m.MongoDB.GetMetadata(),
		eventCollectionName:     defaultEventCollectionName,
		snapshotCollectionName:  defaultSnapshotCollectionName,
		aggregateCollectionName: defaultAggregateCollectionName,
	}
	if val, ok := metadata.Properties[eventCollectionName]; ok && val != "" {
		meta.eventCollectionName = val
	}
	if val, ok := metadata.Properties[snapshotCollectionName]; ok && val != "" {
		meta.snapshotCollectionName = val
	}
	if val, ok := metadata.Properties[aggregateCollectionName]; ok && val != "" {
		meta.aggregateCollectionName = val
	}
	return &meta, nil
}
