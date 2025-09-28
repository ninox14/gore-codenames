package server

import (
	"context"
	"log/slog"
	"math/rand/v2"

	"github.com/ninox14/gore-codenames/internal/database"
	"github.com/ninox14/gore-codenames/internal/database/dto"
	"github.com/ninox14/gore-codenames/internal/database/sqlc"
)

const (
	DefaultAssassinCount   int = 1
	DefaultMaxWordsPerTeam int = 9
)

func GetDefaultBoardSize() *dto.BoardSize {
	return &dto.BoardSize{
		X: 5,
		Y: 5,
	}
}

func CreateEmptyTeam() *dto.Team {
	return &dto.Team{CaptainID: nil, Players: make([]dto.GameStatePlayer, 0), Clues: make([]*dto.Clue, 0)}
}

func GetInitialGameState(ctx context.Context, user *sqlc.User, db *database.DB, logger *slog.Logger) (*dto.GameState, error) {
	wordpack, err := db.Queries.GetWordpack(ctx, 1)
	if err != nil {
		return nil, err
	}

	board := InitBoardStateFromWordPack(wordpack, DefaultMaxWordsPerTeam, DefaultAssassinCount, GetDefaultBoardSize())

	teams := make(map[dto.TeamColor]*dto.Team)
	// TODO: make it less bad
	teams[dto.TeamColorBlue] = CreateEmptyTeam()
	teams[dto.TeamColorRed] = CreateEmptyTeam()
	return &dto.GameState{
		HostID:     user.ID,
		WordPackID: wordpack.ID,
		Spectators: []dto.GameStatePlayer{},
		Teams:      teams,
		Board:      board,
	}, nil
}

func InitBoardStateFromWordPack(wp sqlc.Wordpack, maxWordsPerteam int, maxAssasins int, size *dto.BoardSize) *dto.Board {
	boardSize := size.X * size.Y
	words := SampleArray(wp.Words, boardSize)

	var turnOrder []dto.TeamColor
	if rand.IntN(2) > 0 {
		turnOrder = []dto.TeamColor{dto.TeamColorBlue, dto.TeamColorRed}
	} else {
		turnOrder = []dto.TeamColor{dto.TeamColorRed, dto.TeamColorBlue}
	}

	indxs := Range(boardSize)
	rand.Shuffle(len(indxs), func(i, j int) {
		indxs[i], indxs[j] = indxs[j], indxs[i]
	})
	assassinsIdxs, indxs := Cut(indxs, 0, 1)
	firstTeamIdxs, indxs := Cut(indxs, 0, maxWordsPerteam)
	secondTeamIdxs, _ := Cut(indxs, 0, maxWordsPerteam-1)

	wordsByTeam := make(map[dto.TeamColor][]int)

	wordsByTeam[turnOrder[0]] = firstTeamIdxs
	wordsByTeam[turnOrder[1]] = secondTeamIdxs

	return &dto.Board{
		Size:            size,
		CurrentBoard:    words,
		TurnOrder:       turnOrder,
		MaxWordsPerTeam: maxWordsPerteam,
		AssassinIndexs:  assassinsIdxs,
		GuessedIndexs:   make([]int, 0),
		WordsByTeam:     wordsByTeam,
	}
}

func SampleArray[T any](arr []T, max int) []T {
	set := make(map[int]bool)

	for len(set) < max {
		rand := RandRange(0, len(arr)-1)

		ok := set[rand]
		if !ok {
			set[rand] = true
		}
	}
	var res []T
	for idx := range set {
		res = append(res, arr[idx])
	}
	return res
}

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func Range(n int) []int {
	nums := make([]int, n+1)
	for i := 0; i <= n; i++ {
		nums[i] = i
	}
	return nums
}

// Cut returns the removed elements and the new slice
func Cut[T any](s []T, start, end int) ([]T, []T) {
	if start < 0 || end > len(s) || start > end {
		panic("invalid slice indices")
	}
	cut := append([]T{}, s[start:end]...) // elements cut out
	rest := append(s[:start], s[end:]...) // remaining elements
	return cut, rest
}
