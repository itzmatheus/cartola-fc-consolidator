package repository

import (
	"database/sql"

	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/db"
)

type Repository struct {
	dbConn *sql.DB
	*db.Queries
}
