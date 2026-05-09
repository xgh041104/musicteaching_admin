package main

import (
	"ai_summary_project/cmd/server/wire"
	"ai_summary_project/pkg/config"
	"ai_summary_project/pkg/log"
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	_ "gorm.io/driver/mysql"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
