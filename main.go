package main

import (
	"fmt"
	"github.com/godbus/dbus"
	"os"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}

	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='/org/xmonad/Log',interface='org.xmonad.Log',member='Update'")

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for s := range c {
		if len(s.Body) == 0 {
			continue
		}
		// Convert the result to a slice of strings
		res, ok := s.Body[0].(string)
		if !ok {
			continue
		}
		fmt.Println(res)
	}
}
