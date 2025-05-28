package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "practice"
	collectionName = "gonews"
)

type Storage struct {
	db *mongo.Client
}

func New(uri string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI(uri)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Disconnect(context.Background())

	err = db.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	s := Storage{
		db: db,
	}
	return &s, nil
}

func (s *Storage) Posts() ([]storage.Post, error) {
	db := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	cur, err := db.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	var posts []storage.Post
	for cur.Next(context.Background()) {
		var post storage.Post
		err := cur.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, cur.Err()
}

func (s *Storage) AddPost(post storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdatePost(post storage.Post) error {
	db := s.db.Database(databaseName).Collection(collectionName)

	filter := bson.M{"id": post.ID}
	update := bson.M{
		"$set": bson.M{
			"title":       post.Title,
			"content":     post.Content,
			"authorid":    post.AuthorID,
			"createdat":   post.CreatedAt,
			"publishedat": post.PublishedAt,
		},
	}
	_, err := db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeletePost(post storage.Post) error {
	db := s.db.Database(databaseName).Collection(collectionName)

	filter := bson.M{"id": post.ID}
	_, err := db.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return nil
}
