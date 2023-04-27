package spotify

import (
	"context"
	"fmt"
	"time"

	"github.com/zmb3/spotify/v2"
)

type Spotify struct {
	client *spotify.Client
	ctx    context.Context
}

func NewSpotify(client *spotify.Client) *Spotify {
	return &Spotify{
		client: client,
		ctx: context.Background(),
	}
}

func (s *Spotify) GetPlaylists() ([]spotify.SimplePlaylist, error) {
	playlists, err := s.client.CurrentUsersPlaylists(s.ctx)
	if err != nil {
		return nil, err
	}

	return playlists.Playlists, nil
}

func (s *Spotify) GetDevices() ([]spotify.PlayerDevice, error) {
	devices, err := s.client.PlayerDevices(s.ctx)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (s *Spotify) PlayAlarm(playlistID string, deviceID string, transitionTime float64) error {
	device := spotify.ID(deviceID)
	playlist := spotify.URI("spotify:playlist:" + playlistID)

	fmt.Println("Transferring playback to device: ", device)
	err := s.client.TransferPlayback(s.ctx, device, false)
	if err != nil {
		fmt.Println("Error transferring playback: ", err)
		return err
	}

	// set volume to 0
	fmt.Println("Setting volume to 0")
	err1 := s.client.VolumeOpt(s.ctx, 0, &spotify.PlayOptions{ DeviceID : &device })
	if err1 != nil {
		fmt.Println("Error setting volume to 0: ", err1)
		err2 := s.client.PlayOpt(s.ctx, &spotify.PlayOptions{ PlaybackContext : &playlist, DeviceID : &device })
		if err2 != nil {
			fmt.Println("Error playing playlist: ", err2)
			return err2
		}
		return err1
	}
	
	fmt.Println("Playing playlist: ", playlist)
	err2 := s.client.PlayOpt(s.ctx, &spotify.PlayOptions{ PlaybackContext : &playlist, DeviceID : &device })
	if err2 != nil {
		fmt.Println("Error playing playlist: ", err2)
		return err2
	}

	fmt.Println("Fading volume up to 100 over ", transitionTime, " seconds")
	go s.fade(transitionTime, device)
	return nil
}

func (s *Spotify) fade(transitionTime float64, deviceId spotify.ID) error {
	// fade volume up to 100 over transitionTime 
	// send a request every second
	
	requests := int(transitionTime)
	perRequest := float64(100) / float64(requests)
	fmt.Println("step: ", perRequest)
	for i := 0; i < requests; i++ {
		percent := int(perRequest * (float64(i) + 1))
		fmt.Println("Setting volume to ", percent)
		err := s.client.VolumeOpt(s.ctx, percent, &spotify.PlayOptions{ DeviceID : &deviceId })
		if err != nil {
			fmt.Println("Error setting volume to ", percent, ": ", err)
			return err
		}
		time.Sleep(time.Second)
	}		
	return nil
}
