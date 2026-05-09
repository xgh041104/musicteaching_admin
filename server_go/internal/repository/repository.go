package repository

import (
	"ai_summary_project/pkg/log"
	"ai_summary_project/pkg/zapgorm2"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"

	_ "gorm.io/driver/mysql"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db *gorm.DB
	//rdb    *redis.Client
	//mongo  *mongo.Client
	logger *log.Logger
}

func NewRepository(
	logger *log.Logger,
	db *gorm.DB,
	// rdb *redis.Client,
	//
	//	mongo *mongo.Client,
) *Repository {
	return &Repository{
		db: db,
		//rdb:    rdb,
		//mongo:  mongo,
		logger: logger,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	// ✅ 从配置中获取数据库驱动和 DSN
	driver := conf.GetString("data.db.driver")
	dsn := conf.GetString("data.db.dsn")

	fmt.Println("Connecting to database driver:", driver)
	fmt.Println("DSN:", dsn)

	// GORM 日志配置
	logger := zapgorm2.New(l.Logger)

	// ✅ 根据 driver 建立数据库连接
	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger,
		})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger,
		})
	default:
		panic("unknown db driver: " + driver)
	}

	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// ✅ 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get db instance: %v", err))
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}
func NewMongo(conf *viper.Viper) (*mongo.Client, func(), error) {
	// https://www.mongodb.com/zh-cn/docs/drivers/go/current/
	uri := conf.GetString("data.mongo.uri")
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(fmt.Sprintf("mongo client error: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("mongo ping error: %s", err.Error()))
	}

	return client, func() {
		err = client.Disconnect(ctx)
		if err != nil {
			panic(fmt.Sprintf("mongo disconnect error: %s", err.Error()))
		}
	}, err
}
