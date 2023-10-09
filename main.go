package main

import (
    "net/http"

    "github.com/gin-gonic/gin"

	"github.com/google/uuid"

	"mime/multipart"

	"encoding/json"

	"strings"

	"fmt"
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
var albums = []album{}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Image is required", "error": err.Error()})
		return
	}

	profileData, _ := c.GetPostForm("profile")

	// Clean the profileData before parsing it to JSON
	profileData = cleanProfileString(profileData)

	fmt.Println("Cleaned Profile Data:", profileData)

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

    albums = append(albums, newAlbum)
    c.JSON(http.StatusCreated, gin.H{"albumID": newAlbum.ID,"imageSize": "size"})
}

// This function is used to convert AlbumInfo object passed in as profile data into json format
func cleanProfileString(profileContent string) string {
    profileContent = strings.Replace(profileContent, "class AlbumsProfile {", "{", -1)
	profileContent = strings.Replace(profileContent, "artist: ", "\"artist\": \"", -1)
	profileContent = strings.Replace(profileContent, "title: ", "\",\"title\": \"", -1)
	profileContent = strings.Replace(profileContent, "year: ", "\",\"year\": \"", -1)
	profileContent = strings.Replace(profileContent, "}", "\"}", -1)
	
	// Additional cleanup
	profileContent = strings.Replace(profileContent, "\n    ", "", -1)  // Remove newline and spaces
	profileContent = strings.Replace(profileContent, "\n", "", -1)  // Remove any remaining newline

	return profileContent
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")

    // Loop over the list of albums, looking for
    // an album whose ID value matches the parameter.
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

