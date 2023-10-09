package main

import (
    "net/http"

    "github.com/gin-gonic/gin"

	"github.com/google/uuid"

	"mime/multipart"

	"encoding/json"
)

type profile struct {
	Artist string `json:"artist" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Year   string `json:"year" binding:"required"`
}

type album struct {
	Image  *multipart.FileHeader `form:"image" binding:"required"`
	ID     string  
	Profile profile `json:"profile" binding:"required"`
}

// albums slice to seed record album data.
var albums = []album{
	{
        ID: "1",
        Image: nil,  
        Profile: profile{
            Artist: "artistName",
            Title:  "titleName",
            Year:   "yearValue",
        },
    },
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Image is required", "error": err.Error()})
		return
	}

	profileData, _ := c.GetPostForm("profile")
	var newProfile profile
	err = json.Unmarshal([]byte(profileData), &newProfile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid profile data provided", "error": err.Error()})
		return
	}

	var newAlbum album

    newAlbum.ID = uuid.New().String()
	newAlbum.Profile = newProfile
	newAlbum.Image = image

    // albums = append(albums, newAlbum) // not going to add album to prevent memory overload
    c.JSON(http.StatusCreated, gin.H{"albumID": newAlbum.ID,"imageSize": image.Size})
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")
	id = "1"

    // Loop over the list of albums, looking for
    // an album whose ID value matches 1, which is id of our predefined data
    for _, a := range albums {
        if a.ID == id {
			c.JSON(http.StatusCreated, gin.H{"artist": a.Profile.Artist, "title": a.Profile.Title, "year": a.Profile.Year})
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func main() {
    router := gin.Default()
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

    router.Run(":8080")
}

