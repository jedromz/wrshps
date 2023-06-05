package api

import (
	"fmt"
	"sort"
)

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

// Function to print PlayerStats in a formatted manner
func (ps PlayerStats) String() string {
	return fmt.Sprintf("Nick: %s, Games: %d, Points: %d, Rank: %d, Wins: %d",
		ps.Nick, ps.Games, ps.Points, ps.Rank, ps.Wins)
}

// Function to print TopPlayerStats in a formatted manner
func (tps TopPlayerStats) String() string {
	var result string
	for i, ps := range tps.Stats {
		result += fmt.Sprintf("%d.%s\n", i+1, ps)
	}
	return result
}

// TopPlayerStats represents statistics of top 10 players
type TopPlayerStats struct {
	Stats []PlayerStats `json:"stats"`
}

type Shot struct {
	Coord string `json:"coord"`
}

type GameState struct {
	PlayerBoard  [10][10]string `json:"player_board"`
	OppBoard     [10][10]string `json:"opp_board"`
	TotalShots   int            `json:"total_shots"`
	TotalHits    int            `json:"total_hits"`
	PlayerDesc   string         `json:"player_desc"`
	OppDesc      string         `json:"opp_desc"`
	OppShipsSunk map[int]int
}
type GameStat struct {
	Games  int    `json:"games"`
	Nick   string `json:"nick"`
	Points int    `json:"points"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
}

type GameStats []GameStat

func (g GameStats) Len() int {
	return len(g)
}

func (g GameStats) Less(i, j int) bool {
	if g[i].Wins != g[j].Wins {
		return g[i].Wins > g[j].Wins
	}
	return g[i].Nick < g[j].Nick
}

func (g GameStats) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}
func (g GameStats) String() string {
	// Sort the game statistics before printing
	sort.Sort(g)

	// Format the game statistics as a string
	var result string
	for _, stat := range g {
		result += fmt.Sprintf("Nick: %s\n", stat.Nick)
		result += fmt.Sprintf("Games: %d\n", stat.Games)
		result += fmt.Sprintf("Points: %d\n", stat.Points)
		result += fmt.Sprintf("Rank: %d\n", stat.Rank)
		result += fmt.Sprintf("Wins: %d\n", stat.Wins)
		result += "\n"
	}

	return result
}
