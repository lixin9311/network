package main

import (
	"fmt"
	"net"

	"github.com/lixin9311/network"
)

func main() {
	ip := net.ParseIP("8.8.8.8")
	r, err := network.GetRoute(ip)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("route:", r)
}
