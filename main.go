package main

import (
	"fmt"
	"github.com/Bention99/pokedexcli/internal/config"
)

func main() {
	content := Read()
	fmt.Printf("%v", content)
}