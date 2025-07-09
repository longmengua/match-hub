package mongoclient

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Config struct {
	Hosts                            []string
	Username                         string
	Password                         string
	DatabaseName                     string
	MinPoolSize                      uint64
	MaxPoolSize                      uint64
	MaxConnIdleTime                  time.Duration
	EnableStandardReadWriteSplitMode bool
	MaxStaleness                     time.Duration
	ReplicaName                      string
	Compressors                      []string
}

type MongoClient struct {
	serverConf  *Config
	options     *options.ClientOptions
	db          *mongo.Database
	client      *mongo.Client
	collections map[string]*mongo.Collection
	ctx         context.Context
	cancel      context.CancelFunc
}

func New(conf *Config, opts ...func(*MongoClient)) *MongoClient {
	cli := &MongoClient{serverConf: conf}

	cli.options = options.Client().
		SetHosts(conf.Hosts).
		SetMaxConnIdleTime(30 * time.Second).
		SetMinPoolSize(5).
		SetMaxPoolSize(100)

	if conf.MinPoolSize != 0 {
		cli.options.SetMinPoolSize(conf.MinPoolSize)
	}
	if conf.MaxPoolSize != 0 {
		cli.options.SetMaxPoolSize(conf.MaxPoolSize)
	}
	if conf.MaxConnIdleTime != 0 {
		cli.options.SetMaxConnIdleTime(conf.MaxConnIdleTime)
	}

	if conf.Username != "" {
		cli.options.SetAuth(options.Credential{
			Username: conf.Username,
			Password: conf.Password,
		})
	}

	if conf.EnableStandardReadWriteSplitMode {
		maxStaleness := conf.MaxStaleness
		if maxStaleness < 90*time.Second {
			maxStaleness = 90 * time.Second
		}
		cli.options.SetWriteConcern(writeconcern.Majority())
		cli.options.SetReadPreference(readpref.SecondaryPreferred(readpref.WithMaxStaleness(maxStaleness)))
		cli.options.SetReadConcern(readconcern.Majority())
	}

	if conf.ReplicaName != "" {
		cli.options.SetReplicaSet(conf.ReplicaName)
	}

	if len(conf.Compressors) > 0 {
		cli.options.SetCompressors(conf.Compressors)
	}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

// Start establishes connection, stores db, collection, context
func (cli *MongoClient) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, cli.options)
	if err != nil {
		cancel()
		return fmt.Errorf("mongo connect error: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		return fmt.Errorf("mongo ping error: %w", err)
	}

	db := client.Database(cli.serverConf.DatabaseName)

	cli.ctx = ctx
	cli.cancel = cancel
	cli.client = client
	cli.db = db

	return nil
}

func (m *MongoClient) Close() {
	if m.cancel != nil {
		m.cancel()
	}
	if m.client != nil {
		_ = m.client.Disconnect(context.Background())
	}
}

func (m *MongoClient) Collection(collectionName string) *mongo.Collection {
	if m.collections == nil {
		m.collections = make(map[string]*mongo.Collection)
	}
	if m.collections[collectionName] != nil {
		return m.collections[collectionName]
	}
	m.collections[collectionName] = m.db.Collection(collectionName)
	if m.collections[collectionName] == nil {
		fmt.Printf("Collection %s not found in database %s\n", collectionName, m.serverConf.DatabaseName)
		return nil
	}

	return m.collections[collectionName]
}
