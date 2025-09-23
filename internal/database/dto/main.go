package dto

import "github.com/google/uuid"

type TeamColor string

const (
	TeamColorRed  TeamColor = "red"
	TeamColorBlue TeamColor = "blue"
)

type GameStatePlayer struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Clue struct {
	Word   string `json:"word"`
	Number int    `json:"number"`
}

type Team struct {
	CaptainID *uuid.UUID        `json:"captain_id"`
	Players   []GameStatePlayer `json:"players"` // maybe should be slice of pointers?
	Clues     []*Clue           `json:"clues"`
}

type BoardSize struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Board struct {
	Size            *BoardSize          `json:"size"`
	CurrentBoard    []string            `json:"current_board"`
	GuessedIndexs   []int               `json:"guessed_indexes"`
	AssassinIndexs  []int               `json:"assassin_indexes"`
	TurnOrder       []TeamColor         `json:"turn_order"`
	MaxWordsPerTeam int                 `json:"max_words_per_team"`
	WordsByTeam     map[TeamColor][]int `json:"words_by_team"`
}

type GameState struct {
	HostID     uuid.UUID           `json:"host_id"`
	WordPackID int32               `json:"wordpack_id"`
	Spectators []GameStatePlayer   `json:"spectators"`
	Teams      map[TeamColor]*Team `json:"teams"`
	Board      *Board              `json:"board"`
}
