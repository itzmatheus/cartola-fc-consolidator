package entity

import (
	"strconv"
	"time"
)

type MatchResult struct {
	teamAScore int
	teamBScore int
}

type Match struct {
	ID      string
	TeamA   *Team
	TeamB   *Team
	TeamAID string
	TeamBID string
	Date    time.Time
	Status  string
	Result  MatchResult
	Actions []GameAction
}

func NewMatchResult(teamAScore, teamBScore int) *MatchResult {
	return &MatchResult{
		teamAScore: teamAScore,
		teamBScore: teamBScore,
	}
}

func NewMatch(id string, teamA *Team, teamB *Team, date time.Time) *Match {
	return &Match{
		ID:      id,
		TeamA:   teamA,
		TeamB:   teamB,
		TeamAID: teamA.ID,
		TeamBID: teamB.ID,
		Date:    date,
		Status:  "not started",
	}
}

func (m *MatchResult) GetResult() string {
	return strconv.Itoa(m.teamAScore) + "-" + strconv.Itoa(m.teamBScore)
}
