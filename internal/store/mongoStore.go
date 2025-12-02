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

func (m *MongoStore) pagesCol() *mongo.Collection {
	return m.Client.Database(m.DBName).Collection(colNames.Pages)
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
	log.Printf("[INFO] Saving page '%v'", p.Title)
	p.UpdatedAt = time.Now()
	col := m.pagesCol()

	_, err := col.UpdateOne(
		ctx,
		bson.M{"title": p.Title},
		bson.M{"$set": p},
		options.UpdateOne().SetUpsert(true),
	)

	return err
}

func (m *MongoStore) SearchPages(ctx context.Context, q SearchQuery) ([]*Page, error) {
	if q.Limit == 0 || q.Page == 0 {
		return nil, fmt.Errorf("SearchQuery.Limit and SearchQuery.Page cannot be less than 1")
	}

	log.Printf("[INFO] Searching for pages: '%v', Limit - %d, Page %d", q.Search, q.Limit, q.Page)
	col := m.pagesCol()
	opts := options.Find().
		SetLimit(int64(q.Limit)).
		SetSkip(int64(q.Skip()))

	filter := bson.D{}
	if q.Search != "" {
		filter = bson.D{{Key: "title", Value: bson.D{{Key: "$regex", Value: q.Search}}}}
	}

	cur, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, executeQueryError(err)
	}
	defer cur.Close(ctx)

	return decodeObjects[Page](cur, ctx)
}

func (m *MongoStore) LoadPage(ctx context.Context, title string) (*Page, bool, error) {
	log.Printf("[INFO] Loading page '%v'", title)
	col := m.pagesCol()

	var p Page

	err := col.FindOne(ctx, bson.M{"title": title}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, executeQueryError(err)
	}

	return &p, true, nil
}

func (m *MongoStore) LoadPages(ctx context.Context, q OrderQuery) ([]*Page, error) {
	log.Printf("[INFO] Loading pages: Limit - %d; Qeury by field '%v', Desc %t", q.Limit, q.Field, q.Desc)
	col := m.pagesCol()

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

	cur, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, executeQueryError(err)
	}
	defer cur.Close(ctx)

	return decodeObjects[Page](cur, ctx)
}

func decodeObjects[T any](cur *mongo.Cursor, ctx context.Context) ([]*T, error) {
	r := make([]*T, 0)

	if err := cur.All(ctx, &r); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return r, nil
}

func executeQueryError(err error) error {
	err = fmt.Errorf("failed to execute find query:: %w", err)
	log.Printf("[ERROR] %s", err.Error())
	return err
}
