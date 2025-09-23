export type UserResponse = { name: string; id: string; created_at: string };

export type TokenResponse = {
  AuthenticationToken: string;
  AuthenticationTokenExpiry: string;
};

export type NewGameResponse = {
  game_id: string;
};

// --- Core game state types (mirroring Go structs) ---
// TODO: Remove this after usage
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const teamColors = ['red', 'blue'] as const;
type TeamColor = (typeof teamColors)[number];

export interface GameStatePlayer {
  id: string; // uuid
  name: string;
}

export interface Clue {
  word: string;
  number: number;
}

export interface Team {
  captain_id?: string; // uuid | null
  players: GameStatePlayer[];
  clues: Clue[];
}

export interface BoardSize {
  x: number;
  y: number;
}

export interface Board {
  size: BoardSize;
  current_board: string[];
  guessed_indexes: number[];
  assassin_indexes: number[];
  turn_order: TeamColor[];
  max_words_per_team: number;
  words_by_team: Record<TeamColor, number[]>;
}

export interface GameState {
  host_id: string; // uuid
  wordpack_id: number;
  spectators: GameStatePlayer[];
  teams: Record<TeamColor, Team>;
  board: Board;
}
