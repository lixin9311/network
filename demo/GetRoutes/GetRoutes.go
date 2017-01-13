package main

import (
	"fmt"

	"github.com/lixin9311/network"
)

func main() {
	r, _ := network.GetRoutes()
	for _, v := range r {
		fmt.Println("route:", v)
	}
}
