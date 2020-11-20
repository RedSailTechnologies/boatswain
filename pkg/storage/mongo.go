package storage

import (
	"context"
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo uses a mongodb for storage
type Mongo struct {
	conn string
	db   string
}

// NewMongo returns the Storage interface from the mongo driver
func NewMongo(conn, db string) (Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &Mongo{
		conn: conn,
		db:   db,
	}, nil
}

// IDs gets all uuids in a particular collection
func (m *Mongo) IDs(coll string) ([]string, error) {
	results := make([]string, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.conn))
	if err != nil {
		return nil, err
	}

	c := client.Database(m.db).Collection(coll)
	r, err := c.Distinct(ctx, "uuid", bson.D{})
	if err != nil {
		return nil, err
	}

	for _, uuid := range r {
		results = append(results, uuid.(string))
	}
	return results, nil
}

// GetEvents gets all events for a particular uuid in a particular collection
func (m *Mongo) GetEvents(coll, uuid string) ([]*StoredEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.conn))
	if err != nil {
		return nil, err
	}

	c := client.Database(m.db).Collection(coll)
	cur, err := c.Find(ctx, bson.M{"uuid": uuid})
	if err != nil {
		return nil, err
	}

	results := make([]*StoredEvent, 0)
	for cur.Next(ctx) {
		var result StoredEvent
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	return results, nil
}

// GetVersion gets the most recently inserted event version or 0 on not found/error
func (m *Mongo) GetVersion(coll, uuid string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.conn))
	if err != nil {
		return 0
	}

	c := client.Database(m.db).Collection(coll)
	opts := options.FindOne()
	opts.SetSort(bson.D{primitive.E{Key: "version", Value: -1}})
	result := c.FindOne(ctx, bson.M{"uuid": uuid}, opts)
	if result == nil {
		return 0
	}
	var decoded StoredEvent
	err = result.Decode(&decoded)
	if err != nil {
		return 0
	}

	return decoded.Version
}

// StoreEvent stores an event given the db, uuid, and version
func (m *Mongo) StoreEvent(coll, uuid, eventType, eventData string, version int) error {
	err := m.checkCollection(coll)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.conn))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	c := client.Database(m.db).Collection(coll)
	_, err = c.InsertOne(ctx, StoredEvent{
		UUID:      uuid,
		Version:   version,
		Timestamp: ddd.NewTimestamp(),
		Type:      eventType,
		Data:      eventData,
	})

	return err
}

func (m *Mongo) checkCollection(coll string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.conn))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	colls, err := client.Database(m.db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return err
	}

	exists := false
	for _, name := range colls {
		if name == coll {
			exists = true
		}
	}

	if !exists {
		err = client.Database(m.db).CreateCollection(ctx, coll)
		if err != nil {
			return err
		}
		_, err = client.Database(m.db).Collection(coll).Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.M{
				// TODO AdamP - should this be a constant?
				"uuid":    1,
				"version": 1,
			},
			Options: options.Index().SetUnique(true),
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func truePtr() *bool {
	var t bool
	t = true
	return &t
}
