package repository

import (
	"database/sql"
	"server_go/pkg/log"

	"github.com/spf13/viper"
)

type BaseRepository struct {
	db *sql.DB
	//rdb    *redis.Client
	logger *log.Logger
}

func NewBaseRepository(logger *log.Logger, db *sql.DB) *BaseRepository {
	return &BaseRepository{
		db: db,
		//rdb:    rdb,
		logger: logger,
	}
}

func NewBaseDb(conf *viper.Viper, l *log.Logger) *sql.DB {
	// TODO: init db
	db, err := sql.Open("mysql", conf.GetString("data.mysql.user"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		l.Debug(err.Error())
		return nil
	}
	return db
}
