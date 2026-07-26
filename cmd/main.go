package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/justKody/taskboard-go-api/cmd/api"
	"github.com/justKody/taskboard-go-api/config"
	"github.com/justKody/taskboard-go-api/db"
)

func main() {

	fmt.Printf("%s \n", config.Envs.DatabaseUrl)

	conn := db.NewPostgresStorage(config.Envs.DatabaseUrl)
	ctx := context.Background()

	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Something went wrong when pining the db", err)
	}

	server := api.NewApiServer(fmt.Sprintf(":%s", config.Envs.Port), conn)
	server.Run()
}
