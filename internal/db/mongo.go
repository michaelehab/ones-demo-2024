package db

import (
	"context"
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/ones-demo-2024/internal/config"
	"github.com/project-alvarium/ones-demo-2024/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"sync"
)

type MongoProvider struct {
	cfg    config.MongoConfig
	client *mongo.Client
	logger interfaces.Logger
}

func NewMongoProvider(cfg config.MongoConfig, logger interfaces.Logger) (*MongoProvider, error) {
	mp := &MongoProvider{
		cfg:    cfg,
		logger: logger,
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mp.buildConnectionString()))
	if err != nil {
		return &MongoProvider{}, err
	}
	mp.client = client

	return mp, nil
}

func (mp *MongoProvider) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup) bool {
	wg.Add(1)
	go func() { // Graceful shutdown
		defer wg.Done()

		<-ctx.Done()
		mp.logger.Write(slog.LevelInfo, "shutdown received")
		mp.client.Disconnect(ctx)
	}()
	return true
}

func (mp *MongoProvider) Save(ctx context.Context, item models.SampleData) error {
	recMongo := models.MongoFromSampleData(item)
	coll := mp.client.Database(mp.cfg.DbName).Collection(mp.cfg.Collection)
	result, err := coll.InsertOne(ctx, recMongo)
	if err != nil {
		return err
	}
	mp.logger.Write(slog.LevelDebug, fmt.Sprintf("Document inserted with ID: %s\n", result.InsertedID))
	return nil
}

func (mp *MongoProvider) Close(ctx context.Context) error {
	return mp.client.Disconnect(ctx)
}

func (mp *MongoProvider) buildConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%v", mp.cfg.Username, mp.cfg.Password, mp.cfg.Host, mp.cfg.Port)
}
