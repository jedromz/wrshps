package api

type GameStatus struct {
	GameStatus     string   `json:"game_status"`
	LastGameStatus string   `json:"last_game_status"`
	Nick           string   `json:"nick"`
	OppShots       []string `json:"opp_shots"`
	Opponent       string   `json:"opponent"`
	ShouldFire     bool     `json:"should_fire"`
	Timer          int      `json:"timer"`
}
type StartGameData struct {
	Coords     []string `json:"coords"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick"`
	WPBot      bool     `json:"wpbot"`
}

type GameBoard struct {
	Board []string `json:"board"`
}
type FireData struct {
	Coord string `json:"coord"`
}

type FireResult struct {
	Result string `json:"result"`
}
type GameDescription struct {
	Desc     string `json:"desc"`
	Nick     string `json:"nick"`
	OppDesc  string `json:"opp_desc"`
	Opponent string `json:"opponent"`
}

// GameList represents the list of games
type GameList []struct {
	Guest  string `json:"guest"`
	Host   string `json:"host"`
	ID     string `json:"id"`
	Status string `json:"status"`
}
type LobbyPlayer struct {
	GameStatus string `json:"game_status"`
	Nick       string `json:"nick"`
}

// PlayerStats represents statistics of a player
type PlayerStats struct {
	Games  int    `json:"games"`
	Nick   string `json:"nick"`
	Points int    `json:"points"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
}

// TopPlayerStats represents statistics of top 10 players
type TopPlayerStats struct {
	Stats []PlayerStats `json:"stats"`
}

type Shot struct {
	Coord string `json:"coord"`
}

type GameState struct {
	PlayerBoard [10][10]string `json:"player_board"`
	OppBoard    [10][10]string `json:"opp_board"`
	TotalShots  int            `json:"total_shots"`
	TotalHits   int            `json:"total_hits"`
}
