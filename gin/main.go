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

func main() {

	router := gin.Default() //initializes a Gin router

	/* 	With Gin, you can associate a handler with an HTTP method-and-path combination.
	   	In this way, you can separately route requests sent to a single path based
	   	on the method the client is using. */

	router.GET("/albums", getAlbums)  //associates GET HTTP method and /albums endpoint with a handler
	router.POST("/albums", postAlbum) //associates POST HTTP method and /albums endpoint with a handler
	/* associate GET and /albums/:id path with a handler.
	in Gin, the colon preceding an item in the path signifies that the item is a path parameter. */
	router.GET("/albums/:id", getAlbumByID)

	router.Run("localhost:8080") // attaches the router to an http.Server and start server

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

//postAlbum adds an album from JSON received in request body
func postAlbum(c *gin.Context) {

	var newAlbum album

	//binds the request body to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)

}

func getAlbumByID(c *gin.Context) {

	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusFound, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
