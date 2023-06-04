package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	GameURL        = "/game"
	BoardURL       = "/game/board"
	FireURL        = "/game/fire"
	AbandonGameURL = "/game/abandon"
	GameDescURL    = "/game/desc"
	RefreshURL     = "/game/refresh"
)

type Client struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func NewClient(baseURL string, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) GetGameStatus() (GameStatus, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+GameURL, nil)
	if err != nil {
		return GameStatus{}, err
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return GameStatus{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GameStatus{}, err
	}

	switch resp.StatusCode {
	case 200:
		var gameState GameStatus
		if err := json.Unmarshal(body, &gameState); err != nil {
			return GameStatus{}, err
		}
		return gameState, nil
	case 401:
		var errResp UnauthorizedError
		json.Unmarshal(body, &errResp)
		return GameStatus{}, err
	case 403:
		var errResp ForbiddenError
		json.Unmarshal(body, &errResp)
		return GameStatus{}, err
	case 429:
		var errResp RateLimitExceededError
		json.Unmarshal(body, &errResp)
		return GameStatus{}, err
	default:
		all, _ := ioutil.ReadAll(resp.Body)
		return GameStatus{}, errors.New(fmt.Sprintf("unexpected API error %v %v", resp.StatusCode, fmt.Sprintf("%s", all)))
	}
}
func (c *Client) StartGame(nick, desc, targetNick string, coords []string, botGame bool) (string, error) {
	bodyData := map[string]interface{}{
		"coords":      coords,
		"desc":        desc,
		"nick":        nick,
		"target_nick": targetNick,
		"wpbot":       botGame,
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+GameURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Auth-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case 200:
		token := resp.Header.Get("X-Auth-Token")
		c.Token = token

		return token, nil
	case 400:
		var errResp BadRequestError
		json.Unmarshal(body, &errResp)
		return "", errResp
	case 403:
		var errResp BadRequestError
		json.Unmarshal(body, &errResp)
		return "", errResp
	default:
		return "", errors.New("unexpected API error")
	}
}

func (c *Client) GetGameBoard() (*GameBoard, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+BoardURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
		var gameBoard GameBoard
		if err := json.Unmarshal(body, &gameBoard); err != nil {
			return nil, err
		}
		return &gameBoard, nil
	case 401:
		var errResp UnauthorizedError
		json.Unmarshal(body, &errResp)
		return nil, errResp
	case 403:
		var errResp ForbiddenError
		json.Unmarshal(body, &errResp)
		return nil, errResp
	default:
		return nil, errors.New("unexpected API error")
	}
}
func (c *Client) Fire(data FireData) (FireResult, error) {

	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return FireResult{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+FireURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return FireResult{}, err
	}
	req.Header.Set("X-Auth-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return FireResult{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return FireResult{}, err
	}
	switch resp.StatusCode {
	case 200:
		var fireResult FireResult
		if err := json.Unmarshal(body, &fireResult); err != nil {
			return FireResult{}, err
		}
		return fireResult, err
	case 400:
		var errResp BadRequestError
		json.Unmarshal(body, &errResp)
		return FireResult{}, err
	case 401:
		var errResp UnauthorizedError
		json.Unmarshal(body, &errResp)
		return FireResult{}, err
	case 403:
		var errResp ForbiddenError
		json.Unmarshal(body, &errResp)
		return FireResult{}, err
	case 429:
		var errResp RateLimitExceededError
		json.Unmarshal(body, &errResp)
		return FireResult{}, err
	default:
		return FireResult{}, errors.New("unexpected API error")
	}
}
func (c *Client) AbandonGame() error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+AbandonGameURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 200:
		return nil
	case 400:
		var errResp BadRequestError
		json.Unmarshal(body, &errResp)
		return errResp
	case 401:
		var errResp UnauthorizedError
		json.Unmarshal(body, &errResp)
		return errResp
	case 403:
		var errResp ForbiddenError
		json.Unmarshal(body, &errResp)
		return errResp
	default:
		return errors.New("unexpected API error")
	}
}
func (c *Client) GetGameDescription() (GameDescription, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+GameDescURL, nil)
	if err != nil {
		return GameDescription{}, err
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return GameDescription{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GameDescription{}, err
	}

	switch resp.StatusCode {
	case 200:
		var gameDescription GameDescription
		if err := json.Unmarshal(body, &gameDescription); err != nil {
			return GameDescription{}, err
		}
		return gameDescription, nil
	case 401:
		var errResp UnauthorizedError
		json.Unmarshal(body, &errResp)
		return GameDescription{}, err
	case 404:
		var errResp NotFoundError
		json.Unmarshal(body, &errResp)
		return GameDescription{}, err
	case 429:
		var errResp RateLimitExceededError
		json.Unmarshal(body, &errResp)
		return GameDescription{}, err
	default:
		return GameDescription{}, errors.New("unexpected API error")
	}
}
func (c *Client) RefreshGameSession() error {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+RefreshURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 200:
		return nil
	case 400:
		var errResp BadRequestError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return err
		}
		return errResp
	case 401:
		var errResp UnauthorizedError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return err
		}
		return errResp
	case 403:
		var errResp ForbiddenError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return err
		}
		return errResp
	case 429:
		var errResp RateLimitExceededError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return err
		}
		return errResp
	default:
		return errors.New("unexpected API error")
	}
}
func (c *Client) GetAllGames(status string) (GameList, error) {
	var gameList GameList

	req, err := http.NewRequest(http.MethodGet, c.BaseURL+"/list", nil)
	if err != nil {
		return gameList, err
	}

	q := req.URL.Query()
	q.Add("status", status)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return gameList, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return gameList, err
	}

	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(body, &gameList); err != nil {
			return gameList, err
		}
		return gameList, nil
	case 400:
		var errResp BadRequestError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return gameList, err
		}
		return gameList, errResp
	default:
		return gameList, errors.New("unexpected API error")
	}
}

// GetLobbyPlayers retrieves a list of players waiting in the lobby
func (c *Client) GetLobbyPlayers() ([]LobbyPlayer, error) {
	var players []LobbyPlayer

	req, err := http.NewRequest(http.MethodGet, c.BaseURL+"/lobby", nil)
	if err != nil {
		return players, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return players, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return players, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp BadRequestError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return players, err
		}
		return players, errResp
	}

	if err := json.Unmarshal(body, &players); err != nil {
		return players, err
	}

	return players, nil
}

// GetTopPlayerStats retrieves top 10 players' statistics
func (c *Client) GetTopPlayerStats() (TopPlayerStats, error) {
	var topStats TopPlayerStats

	req, err := http.NewRequest(http.MethodGet, c.BaseURL+"/stats", nil)
	if err != nil {
		return topStats, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return topStats, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return topStats, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp BadRequestError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return topStats, err
		}
		return topStats, errResp
	}

	if err := json.Unmarshal(body, &topStats); err != nil {
		return topStats, err
	}

	return topStats, nil
}

func (c *Client) GetPlayerStats(nick string) (GameStats, error) {
	url := "https://go-pjatk-server.fly.dev/api/stats/" + strings.TrimSpace(nick)

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	// Read the response body
	var response struct {
		Stats GameStat `json:"stats"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return GameStats{response.Stats}, nil
}
