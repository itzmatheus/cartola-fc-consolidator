package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/db"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/repository"
	"github.com/itzmatheus/cartola-fc-consolidator/pkg/uow"
	_ "github.com/lib/pq"
)

const (
	DB_HOST = "localhost"
	DB_PORT = 5432
	DB_USER = "consolidator"
	DB_PASS = "consolidator"
	DB_NAME = "consolidator"
)

func main() {

	ctx := context.Background()

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	uow, err := uow.NewUow(ctx, db)
	if err != nil {
		panic(err)
	}
	registerRepositories(uow)

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
