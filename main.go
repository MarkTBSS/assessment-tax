package main

import (
	"github.com/MarkTBSS/assessment-tax/config"
	"github.com/MarkTBSS/assessment-tax/databases"
	"github.com/MarkTBSS/assessment-tax/server"
)

func main() {
	conf := config.ConfigGetting()
	db := databases.NewPostgresDatabase(conf.Database)
	server := server.NewEchoServer(conf, db)
	server.Start()
}
