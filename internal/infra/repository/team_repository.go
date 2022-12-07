package repository

import (
	"context"
	"database/sql"

	"github.com/itzmatheus/cartola-fc-consolidator/internal/domain/entity"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/db"
)

type TeamRepository struct {
	Repository
}

func NewTeamRepository(dbConn *sql.DB) *TeamRepository {
	return &TeamRepository{
		Repository: Repository{
			dbConn:  dbConn,
			Queries: db.New(dbConn),
		},
	}
}

func (r *TeamRepository) AddScore(ctx context.Context, player *entity.Player, score float64) error {
	err := r.Queries.AddScoreToTeam(ctx, db.AddScoreToTeamParams{
		ID:    player.ID,
		Score: score,
	})
	return err
}

func (r *TeamRepository) FindByID(ctx context.Context, id string) (*entity.Team, error) {
	team, err := r.Queries.FindTeamById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		ID:   team.ID,
		Name: team.Name,
	}, nil
}
