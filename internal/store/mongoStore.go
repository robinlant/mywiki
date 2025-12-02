package store

import (
	"context"
	"errors"

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
	col := m.Client.Database(m.DBName).Collection(colNames.Pages)

	var p Page

	err := col.FindOne(ctx, bson.M{"title": title}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &p, true, nil
}
