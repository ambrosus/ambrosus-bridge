package main

import (
	"relay/config"
	"relay/networks/amb"
)

func main() {
	ambNetwork := config.Networks["amb"]

	for _, bridge := range config.Bridges {
		ambB := amb.New(bridge.Amb)
		sideB := side.New(bridge.Side, bridge.SideNetwork)
		side.
			ambListener := amb.NewListener(bridge.Amb)

	}

}
