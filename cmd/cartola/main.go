package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/db"
	httpHandler "github.com/itzmatheus/cartola-fc-consolidator/internal/infra/http"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/repository"
	"github.com/itzmatheus/cartola-fc-consolidator/pkg/uow"
	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "consolidator"
	DB_PASS     = "consolidator"
	DB_NAME     = "consolidator"
	SERVER_PORT = "8080"
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

	log.Println("Server running at port: " + SERVER_PORT)
	if err = http.ListenAndServe(":"+SERVER_PORT, r); err != nil {
		log.Println("Not working")
		panic(err)
	}

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
