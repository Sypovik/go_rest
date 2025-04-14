package main

import (
	"go-restapi/internal/handler"
	"go-restapi/internal/repository"
	"go-restapi/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, _ := repository.SetupDatabase()
	var repo repository.NoteRepository = repository.NewNoteRepository(db)
	var serv service.NoteService = service.NewNoteService(repo)
	var hand handler.NoteHandler = handler.NewNoteHandler(serv)
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	hand.RegisterRoutes(r)

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // No Content
	})

	r.Run(":8083")
}
