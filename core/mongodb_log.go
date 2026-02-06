package core

import (
	"encoding/json"
	"order_food/global"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

type MongoDBWriteSyncer struct {
	client     *mgo.Database
	collection string
	mu         sync.Mutex
}

func NewMongoDBWriteSyncer(collectionName string) (*MongoDBWriteSyncer, error) {

	return &MongoDBWriteSyncer{
		client:     global.GVA_MONGO,
		collection: collectionName,
	}, nil
}

func (m *MongoDBWriteSyncer) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	doc := make(map[string]interface{})
	json.Unmarshal(p, &doc)
	doc["timestamp"] = time.Now()

	err = m.client.C(m.collection).Insert(&doc)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (m *MongoDBWriteSyncer) Sync() error {
	return nil
}
