package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var mutex = &sync.RWMutex{}

func main() {
	r := gin.Default()
	// this may cause a
	// fatal error: concurrent map iteration and map write
	var hubs = make(map[*Hub]bool)

	r.GET("/:id", func(c *gin.Context) {
		fmt.Println(len(hubs))
		// Adds multi-hub functionality to:
		// https://yourbasic.org/golang/append-explained/
		id := c.Param("id")
		var currentHub *Hub
		// TODO: How do we destroy hubs when empty?
		// Can we set up a defer function somewhere that
		// returns out of the goroutine?
		mutex.Lock()
		for hub := range hubs {
			if hub.Id == id {
				currentHub = hub
			}
		}
		mutex.Unlock()
		if currentHub == nil {
			currentHub = newHub(id)
			fmt.Println(currentHub)
			mutex.Lock()
			hubs[currentHub] = true
			mutex.Unlock()
			go currentHub.run()
		}
		serveWs(currentHub, c.Writer, c.Request)
	})
	r.Run()
}
