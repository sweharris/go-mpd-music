package main

import (
	"os"

	"github.com/fhs/gompd/v2/mpd"
)

// Global for simplicity
var conn *mpd.Client

func get_mpd_addr() (string, string) {
	// If MPD_HOST isn't defined, use the local socket
	addr := os.Getenv("MPD_HOST")
	if addr == "" {
		return "unix", "/var/run/mpd/socket"
	}

	// If it begines with a / then assume it's a local socket
	if addr[0] == '/' {
		return "unix", addr
	}

	port := os.Getenv("MPD_PORT")
	if len(port) == 0 {
		port = "6600"
	}
	return "tcp", addr + ":" + port
}

// Connect to MPD server
func connect_to_mpd() {
	if conn == nil {
		proto, addr := get_mpd_addr()
		c, err := mpd.Dial(proto, addr)
		if err != nil {
			die(err)
		}
		conn = c
	}
}

func get_status() mpd.Attrs {
	connect_to_mpd()
	status, err := conn.Status()
	if err != nil {
		die(err)
	}

	return status
}

func get_song() mpd.Attrs {
	connect_to_mpd()
	song, err := conn.CurrentSong()
	if err != nil {
		die(err)
	}

	return song
}

func load_playlist(list string) {
	connect_to_mpd()
	conn.Clear()
	err := conn.PlaylistLoad(list, -1, -1)
	if err != nil {
		die(err)
	}
}

func get_playlist(s, e int) []mpd.Attrs {
	connect_to_mpd()

	list, err := conn.PlaylistInfo(s, e)
	if err != nil {
		die(err)
	}
	return list
}

func mpd_db_info() []mpd.Attrs {
	connect_to_mpd()

	list, err := conn.ListAllInfo("/")
	if err != nil {
		die(err)
	}
	return list
}
