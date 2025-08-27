package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// Brakująca definicja typu
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
	Models *mongo.Collection
}

type ModelInfo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Provider string             `bson:"provider" json:"provider"`
}

// InitMongoDB tworzy klienta i wybiera bazę oraz kolekcję models
func InitMongoDB(mongoURI string, dbName string) (*MongoDB, error) {
	if mongoURI == "" {
		return nil, errors.New("empty mongo uri")
	}
	if dbName == "" {
		return nil, errors.New("empty db name")
	}

	clientOpts := options.Client().ApplyURI(mongoURI)

	ctxConn, cancelConn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConn()

	client, err := mongo.Connect(ctxConn, clientOpts)
	if err != nil {
		return nil, err
	}

	ctxPing, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()
	if err := client.Ping(ctxPing, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}

	db := client.Database(dbName)
	m := &MongoDB{
		Client: client,
		DB:     db,
		Models: db.Collection("models"),
	}

	if err := m.ensureIndexes(); err != nil {
		log.Printf("warn ensureIndexes failed: %v", err)
	}

	log.Println("Connected to MongoDB and selected DB:", dbName)
	return m, nil
}

func (db *MongoDB) ensureIndexes() error {
	if db == nil || db.Models == nil {
		return errors.New("nil db or collection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Indeks złożony po name i provider z unikalnością
	compound := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: 1}, {Key: "provider", Value: 1}},
		Options: options.Index().
			SetBackground(true).
			SetName("uniq_name_provider").
			SetUnique(true),
	}
	// Dodatkowy prosty po name dla sortowania
	byName := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: 1}},
		Options: options.Index().
			SetBackground(true).
			SetName("idx_name_asc"),
	}

	if _, err := db.Models.Indexes().CreateOne(ctx, compound); err != nil {
		return err
	}
	if _, err := db.Models.Indexes().CreateOne(ctx, byName); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) Close() {
	if db == nil || db.Client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB.")
	}
}

// Ręczny seed
func (db *MongoDB) InsertModels(models []ModelInfo) error {
	if db == nil || db.Models == nil {
		return errors.New("db not initialized")
	}
	if len(models) == 0 {
		return nil
	}

	docs := make([]interface{}, 0, len(models))
	for _, model := range models {
		docs = append(docs, model)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	_, err := db.Models.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	log.Printf("Inserted %d models into MongoDB.", len(models))
	return nil
}

func (db *MongoDB) GetModelByID(idHex string) (*ModelInfo, error) {
	if db == nil || db.Models == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if idHex == "" {
		return nil, fmt.Errorf("empty id")
	}
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, fmt.Errorf("invalid object id: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOne().
		SetProjection(bson.M{"_id": 1, "name": 1, "provider": 1})

	var out ModelInfo
	err = db.Models.FindOne(ctx, bson.M{"_id": oid}, opts).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Jedna strona wyników
func (db *MongoDB) GetModels(page, limit int64) ([]ModelInfo, error) {
	if db == nil || db.Models == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().
		SetProjection(bson.M{"_id": 1, "name": 1, "provider": 1}).
		SetSort(bson.D{{Key: "name", Value: 1}, {Key: "_id", Value: 1}}).
		SetSkip((page - 1) * limit).
		SetLimit(limit)

	cur, err := db.Models.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]ModelInfo, 0, limit)
	for cur.Next(ctx) {
		var m ModelInfo
		if err := cur.Decode(&m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, cur.Err()
}
func (db *MongoDB) GetModelByName(name string) (*ModelInfo, error) {
	if db == nil || db.Models == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if name == "" {
		return nil, fmt.Errorf("empty name")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// jeśli masz wiele rekordów o tej samej nazwie, można dodać sort po provider lub _id
	opts := options.FindOne().
		SetProjection(bson.M{"_id": 1, "name": 1, "provider": 1})

	var out ModelInfo
	err := db.Models.FindOne(ctx, bson.M{"name": name}, opts).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}
