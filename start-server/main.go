// https://gist.github.com/enricofoltran/10b4a980cd07cb02836f70a4ab3e72d7 looks cool
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/carlokuiper/k-in-a-go"
)

func main() {
	game := &kinago.Game{}
	config := kinago.Config{M: 3, N: 3, K: 3}
	if m, err := strconv.Atoi(os.Getenv("M")); err == nil && m != 0 {
		config.M = m
	}
	if n, err := strconv.Atoi(os.Getenv("N")); err == nil && n != 0 {
		config.N = n
	}
	if k, err := strconv.Atoi(os.Getenv("K")); err == nil && k != 0 {
		config.K = k
	}
	game.New(config)
	http.HandleFunc("/status", game.Status)
	http.HandleFunc("/start", game.Start)
	http.HandleFunc("/move", game.Move)

	fileServer := http.FileServer(http.Dir("../static"))
	http.Handle("/", fileServer)

	fmt.Println("Starting server at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
