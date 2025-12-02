package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var colNames = struct {
	Pages string
}{Pages: "pages"}

type MongoStore struct {
	Client *mongo.Client
	DBName string
}

type DisconnectMongo func()

func NewMongoStore(uri string, db string) (*MongoStore, DisconnectMongo) {
	o := options.Client().ApplyURI(uri)
	c, err := mongo.Connect(o)

	if err != nil {
		panic(err)
	}

	return &MongoStore{Client: c, DBName: db}, func() {
		if err := c.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func (m *MongoStore) SavePage(ctx context.Context, p *Page) error {
	log.Printf("[INFO] Saving page %v", p.Title)
	p.UpdatedAt = time.Now()
	coll := m.Client.Database(m.DBName).Collection(colNames.Pages)

	_, err := coll.UpdateOne(
		ctx,
		bson.M{"title": p.Title},
		bson.M{"$set": p},
		options.UpdateOne().SetUpsert(true),
	)

	return err
}

func (m *MongoStore) LoadPage(ctx context.Context, title string) (*Page, bool, error) {
	log.Printf("[INFO] Loading page %v", title)
	col := m.Client.Database(m.DBName).Collection(colNames.Pages)

	var p Page

	err := col.FindOne(ctx, bson.M{"title": title}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to decode document: %w", err)
	}

	return &p, true, nil
}

func (m *MongoStore) LoadPages(ctx context.Context, q Query) ([]*Page, error) {
	log.Printf("[INFO] Loading pages: Limit - %d; Qeury by field %v, Desc %t", q.Limit, q.Field, q.Desc)
	col := m.Client.Database(m.DBName).Collection("pages")

	opts := options.Find().
		SetLimit(int64(q.Limit))

	filter := bson.M{}

	var sortBSON bson.D
	sortByField := q.Field

	if sortByField != "" {
		filter = bson.M{sortByField: bson.M{"$ne": nil}}

		var sortDirection int
		if !q.Desc {
			sortDirection = 1
		} else {
			sortDirection = -1
		}
		sortBSON = bson.D{{Key: sortByField, Value: sortDirection}}
	} else {
		sortBSON = bson.D{{Key: "title", Value: 1}}
	}

	opts.SetSort(sortBSON)

	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	var pages []*Page

	if err := cursor.All(ctx, &pages); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return pages, nil
}
