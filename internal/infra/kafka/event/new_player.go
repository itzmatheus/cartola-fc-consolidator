package event

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/usecase"
	"github.com/itzmatheus/cartola-fc-consolidator/pkg/uow"
)

type ProcessNewPlayer struct{}

func (p ProcessNewPlayer) Process(ctx context.Context, msg *kafka.Message, uow uow.UowInterface) error {
	var input usecase.AddPlayerInput
	err := json.Unmarshal(msg.Value, &input)
	if err != nil {
		return err
	}
	addNewPlayerUseCase := usecase.NewAddPlayerUseCase(uow)
	err = addNewPlayerUseCase.Execute(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
