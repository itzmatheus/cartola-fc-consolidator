package usecase

import (
	"context"

	"github.com/itzmatheus/cartola-fc-consolidator/internal/domain/entity"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/domain/repository"
	"github.com/itzmatheus/cartola-fc-consolidator/pkg/uow"
)

type AddPlayerInput struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	InitialPrice float64 `json:"initial_price"`
}

type AddPlayerUseCase struct {
	Uow uow.UowInterface
}

func NewAddPlayerUseCase(uow uow.UowInterface) *AddPlayerUseCase {
	return &AddPlayerUseCase{
		Uow: uow,
	}
}

func (a *AddPlayerUseCase) getPlayerRepository(ctx context.Context) repository.PlayerRepositoryInterface {
	playerRepository, err := a.Uow.GetRepository(ctx, "PlayerRepository")
	if err != nil {
		panic(err)
	}
	return playerRepository.(repository.PlayerRepositoryInterface)
}

func (a *AddPlayerUseCase) Execute(ctx context.Context, input AddPlayerInput) error {
	playerRepository := a.getPlayerRepository(ctx)
	player := entity.NewPlayer(input.ID, input.Name, input.InitialPrice)
	err := playerRepository.Create(ctx, player)
	if err != nil {
		return err
	}
	a.Uow.CommitOrRollback()
	return nil
}
