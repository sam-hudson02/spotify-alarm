package main

import (
	"fmt"
	"strconv"
	"spotify-alarm/src/spotify"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	ch = make(chan *spotify.Spotify)
)

func main() {
	router := gin.Default()
	spotifyUrl := spotify.GetUrl()
	fmt.Println("Go the following URL to authorize the application: ", spotifyUrl)

	token, err := spotify.ParseToken()
	if err != nil {
		fmt.Println("Couldn't parse token: ", err)
	} else {
		go useToken(token)
	}

	router.GET("/", func(g *gin.Context) {
		g.Redirect(302, spotifyUrl)
	})

	router.GET("/callback", callback)
	router.GET("/trigger", trigger)

	router.Run(":8080")
}


func callback(g *gin.Context) {
	r := g.Request
	w := g.Writer

	fmt.Println("Callback called!")

	client, err := spotify.ParseCode(r)

	if err != nil {
		fmt.Fprintf(w, "Couldn't get client: %v", err)
		return
	}

	g.String(200, "Successfully got client!")
	fmt.Println("Successfully got client!")
	s := spotify.NewSpotify(client)
	go setSpotify(s)
}

func setSpotify(s *spotify.Spotify) {
	ch <- s
}


func useToken(token *oauth2.Token) {
	client := spotify.UseToken(token)
	s := spotify.NewSpotify(client)
	ch <- s
}

func trigger(g *gin.Context) {
	fmt.Println("Trigger called!")
	playlistID := g.Query("playlist_id")
	deviceID := g.Query("device_id")
	transition := g.Query("transition_time")
	// parse transition time to float
	transitionTime, err := strconv.ParseFloat(transition, 64)
	if err != nil {
		g.String(500, "Error parsing transition time: %v", err)
		return
	}

	fmt.Println("Playing alarm!")
	s := <-ch
	fmt.Println("Got spotify client!")
	go s.PlayAlarm(playlistID, deviceID, transitionTime)
	ch <- s
}
