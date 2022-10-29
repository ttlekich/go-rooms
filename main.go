package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	var hubs = make(map[*Hub]bool)

	fmt.Println(hubs)
	r.GET("/:id", func(c *gin.Context) {
		// Adds multi-hub functionality to:
		// https://yourbasic.org/golang/append-explained/
		id := c.Param("id")
		var currentHub *Hub
		// TODO: How do we destroy hubs when empty?
		// Can we set up a defer function somewhere that
		// returns out of the goroutine?
		for hub := range hubs {
			if hub.Id == id {
				currentHub = hub
			}
		}
		if currentHub == nil {
			currentHub = newHub(id)
			fmt.Println(currentHub)
			hubs[currentHub] = true
			go currentHub.run()
		}
		serveWs(currentHub, c.Writer, c.Request)
	})
	r.Run()
}
