package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/carlokuiper/k-in-a-go"
)

func main() {
	move := kinago.Coordinate{}
	if x, err := strconv.Atoi(os.Getenv("X")); err == nil && x != 0 {
		move.X = x
	}
	if y, err := strconv.Atoi(os.Getenv("Y")); err == nil && y != 0 {
		move.Y = y
	}
	body, err := json.Marshal(move)
	if err != nil {
		panic(err)
	}
	r, err := http.NewRequest("POST", "http://localhost:8080/turn", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer func() { _ = response.Body.Close() }()
	game := kinago.Game{}
	if err = json.NewDecoder(response.Body).Decode(&game); err != nil {
		panic(err)
	}
	fmt.Println(response.StatusCode, game.Board)
}
