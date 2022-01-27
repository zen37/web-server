package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

/*
getAlbums responds with the list of all albums as JSON.

gin.Context is the most important part of Gin; it carries request details, validates and serializes JSON, and more.
Despite the similar name, this is different from Go’s built-in context package.
*/

func getAlbums(c *gin.Context) {
	//Context.IndentedJSON can be replaced with a call to Context.JSON to send more compact JSON.
	//In practice, the indented form is much easier to work with when debugging and the size difference is usually small.
	c.IndentedJSON(http.StatusOK, albums)
	//c.JSON(http.StatusOK, albums)
}

func main() {

	router := gin.Default()          //initializes a Gin router
	router.GET("/albums", getAlbums) //associates GET HTTP method and /albums endpoint with a handler
	router.Run("localhost:8080")     // attaches the router to an http.Server and start server

}
