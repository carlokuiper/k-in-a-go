package kinago

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Game struct {
	K        int          `json:"k"`
	Board    [][]int      `json:"board"`
	History  []Coordinate `json:"history"`
	Finished []Coordinate `json:"finished"`
	mu       sync.Mutex
}

func (g *Game) New(config Config) {
	board := make([][]int, config.N) // create rows
	for i := range board {
		board[i] = make([]int, config.M) // for each row add the columns
	}
	g.K = config.K
	g.Board = board
	g.History = nil
	g.Finished = nil
}

func (g *Game) Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var config Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()
	if !config.valid() {
		http.Error(w, "invalid config", http.StatusUnprocessableEntity)
		return
	}
	g.mu.Lock()
	g.New(config)
	g.mu.Unlock()
	g.writeResponse(w)
}
func (g *Game) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	response, err := json.Marshal(g)
	if err != nil {
		http.Error(w, "cannot marshall response", http.StatusInternalServerError)
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "cannot write response", http.StatusInternalServerError)
	}
	g.writeResponse(w)
}

func (g *Game) Move(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var move Coordinate
	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()
	g.mu.Lock()
	nextValue := g.nextValue()
	if nextValue == 0 {
		http.Error(w, "cannot determine next value", http.StatusInternalServerError)
	}
	if err := g.update(move, nextValue); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	g.mu.Unlock()
	g.writeResponse(w)
}

func (g *Game) writeResponse(w http.ResponseWriter) {
	response, err := json.Marshal(g)
	if err != nil {
		http.Error(w, "cannot marshall response", http.StatusInternalServerError)
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "cannot write response", http.StatusInternalServerError)
	}
}

func (g *Game) nextValue() int {
	if g.History == nil {
		return 1
	}
	previousMove := g.History[len(g.History)-1]
	switch g.Board[previousMove.Y][previousMove.X] {
	case 1:
		return 2
	case 2:
		return 1
	default:
		return 0
	}
}

func (g *Game) update(move Coordinate, nextValue int) error {
	m, n := mxn(g.Board)
	if move.X < 0 || move.Y < 0 || move.X >= n || move.Y >= m {
		return fmt.Errorf("move invalid")
	}
	if g.History == nil {
		g.Board[move.Y][move.X] = nextValue
		g.History = []Coordinate{move}
		return nil
	}
	previousMove := g.History[len(g.History)-1]
	if previousMove == move {
		return nil
	}
	if g.Finished != nil {
		return fmt.Errorf("game already finished")
	}
	if g.Board[move.Y][move.X] != 0 {
		return fmt.Errorf("move not free")
	}
	g.Board[move.Y][move.X] = nextValue
	g.History = append(g.History, move)
	g.Finished = finished(g.Board, g.K)
	return nil
}

type position struct {
	Coordinate
	value int
}

func finished(board [][]int, k int) []Coordinate {
	m, n := mxn(board)
	for y, r := range board {
		row := make([]position, 0, k)
		for x, value := range r {
			row = append(row, position{Coordinate{x, y}, value})
		}
		if result := kInARow(row, k); result != nil {
			return result
		}
	}
	for x := range m {
		column := make([]position, n)
		for y, row := range board {
			column[y] = position{Coordinate{x, y}, row[x]}
		}
		if result := kInARow(column, k); result != nil {
			return result
		}
	}
	// (off) diagonal
	for i := -(n - 1); i < m; i++ {
		x := i
		y := 0
		offDiagonal := make([]position, 0, n)
		for range n {
			if x >= 0 && y >= 0 && x < m && y < n {
				offDiagonal = append(offDiagonal, position{Coordinate{x, y}, board[y][x]})
			}
			x++
			y++
		}
		if result := kInARow(offDiagonal, k); result != nil {
			return result
		}
	}
	// (off) anti diagonal
	for i := 0; i < m+n; i++ {
		x := 0
		y := i
		offAntiDiagonal := make([]position, 0, m)
		for range m {
			if x >= 0 && y >= 0 && x < m && y < n {
				offAntiDiagonal = append(offAntiDiagonal, position{Coordinate{x, y}, board[y][x]})
			}
			x++
			y--
		}
		if result := kInARow(offAntiDiagonal, k); result != nil {
			return result
		}
	}
	return nil
}

func mxn(board [][]int) (int, int) {
	n := len(board)
	if n == 0 {
		return 0, 0
	}
	return len(board[0]), n
}

func kInARow(row []position, k int) []Coordinate {
	var lastEl int
	var result []Coordinate
	for _, p := range row {
		switch p.value {
		case 0:
			result = nil
			lastEl = 0
		case lastEl:
			result = append(result, p.Coordinate)
			if len(result) == k {
				return result
			}
		default:
			result = []Coordinate{p.Coordinate}
			lastEl = p.value
		}
	}
	return nil
}

type Config struct {
	M int `json:"m"` // board dimensions are m x n
	N int `json:"n"` // board dimensions are m x n
	K int `json:"k"` // k-in-a-row to win
}

func (c *Config) valid() bool {
	if c.M <= 0 || c.N <= 0 || c.K <= 0 {
		return false
	}
	if c.K > c.M && c.K > c.N {
		return false
	}
	return true
}
