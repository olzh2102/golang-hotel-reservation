package db

import (
	"context"
	"os"

	"github.com/olzh2102/golang-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, Map, Map) error
	GetHotels(context.Context, Map, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection("hotels"),
	}
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter, update Map) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter Map, pag *Pagination) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((pag.Page - 1) * pag.Limit)
	opts.SetLimit(pag.Limit)
	res, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := res.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, Map{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}
