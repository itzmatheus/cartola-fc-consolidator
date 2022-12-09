package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/db"
	httpHandler "github.com/itzmatheus/cartola-fc-consolidator/internal/infra/http"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/kafka/consumer"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/repository"
	"github.com/itzmatheus/cartola-fc-consolidator/pkg/uow"
	_ "github.com/lib/pq"
)

const (
	DB_HOST      = "172.17.0.1"
	DB_PORT      = 5432
	DB_USER      = "consolidator"
	DB_PASS      = "consolidator"
	DB_NAME      = "consolidator"
	SERVER_PORT  = "8080"
	KAFKA_BROKER = "broker:9094"
)

var (
	KAFKA_TOPICS = []string{"newMatch", "chooseTeam", "newPlayer", "matchUpdateResult", "newAction"}
)

func main() {

	ctx := context.Background()

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME)

	dtb, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	err = dtb.Ping()
	if err != nil {
		panic(err)
	}
	defer dtb.Close()

	uow, err := uow.NewUow(ctx, dtb)
	if err != nil {
		panic(err)
	}
	registerRepositories(uow)

	r := chi.NewRouter()
	r.Get("/players", httpHandler.ListPlayerHandler(ctx, *db.New(dtb)))
	r.Get("/myTeam/{teamID}/players", httpHandler.ListMyTeamPlayers(ctx, *db.New(dtb)))
	r.Get("/myTeam/{teamID}/balance", httpHandler.GetMyTeamBalanceHandler(ctx, *db.New(dtb)))
	r.Get("/matches", httpHandler.ListMatchesHandler(ctx, repository.NewMatchRepository(dtb)))
	r.Get("/matches/{matchID}", httpHandler.ListMatchByIDHandler(ctx, repository.NewMatchRepository(dtb)))

	go http.ListenAndServe(":"+SERVER_PORT, r)
	log.Println("Server running at port: " + SERVER_PORT)

	msgChan := make(chan *kafka.Message)
	go consumer.Consume(KAFKA_TOPICS, KAFKA_BROKER, msgChan)
	consumer.ProcessEvents(ctx, msgChan, uow)
}

func registerRepositories(uow *uow.Uow) {
	uow.Register("PlayerRepository", func(tx *sql.Tx) interface{} {
		repo := repository.NewPlayerRepository(uow.Db)
		repo.Queries = db.New(tx)
		return repo
	})

	uow.Register("MatchRepository", func(tx *sql.Tx) interface{} {
		repo := repository.NewMatchRepository(uow.Db)
		repo.Queries = db.New(tx)
		return repo
	})

	uow.Register("TeamRepository", func(tx *sql.Tx) interface{} {
		repo := repository.NewTeamRepository(uow.Db)
		repo.Queries = db.New(tx)
		return repo
	})

	uow.Register("MyTeamRepository", func(tx *sql.Tx) interface{} {
		repo := repository.NewMyTeamRepository(uow.Db)
		repo.Queries = db.New(tx)
		return repo
	})
}
