package service

import (
	"errors"

	"github.com/itzmatheus/cartola-fc-consolidator/internal/domain/entity"
)

func ChoosePlayers(myTeam entity.MyTeam, players []entity.Player) error {
	totalCost := 0.0
	totalEarned := 0.0

	for _, player := range players {

		// Sell player
		if playerInMyTeam(player, myTeam) && !playerInPlayersList(player, &players) {
			totalEarned += player.Price
		}
		// Buy player
		if !playerInMyTeam(player, myTeam) && playerInPlayersList(player, &players) {
			totalCost += player.Price
		}
	}

	totalAvailableForChoosePlayers := myTeam.Score + totalEarned

	if totalCost > totalAvailableForChoosePlayers {
		return errors.New("not enough money")
	}

	myNewTotalScoreAfterChoosePlayers := myTeam.Score + totalEarned - totalCost

	myTeam.Score = myNewTotalScoreAfterChoosePlayers
	myTeam.Players = []string{}

	for _, player := range players {
		myTeam.Players = append(myTeam.Players, player.ID)
	}
	return nil
}

func playerInMyTeam(player entity.Player, myTeam entity.MyTeam) bool {
	for _, playerID := range myTeam.Players {
		if player.ID == playerID {
			return true
		}
	}
	return false
}

func playerInPlayersList(player entity.Player, players *[]entity.Player) bool {
	for _, p := range *players {
		if player.ID == p.ID {
			return true
		}
	}
	return false
}
