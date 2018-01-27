package main

import (
	"fmt"
	"github.com/prestonTao/upnp"
)

func main() {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(55789, 55789, "TCP"); err == nil {
		fmt.Println("success !")
		// remove port mapping in gatway
		mapping.Reclaim()
	} else {
		fmt.Println("fail !")
	}

}
