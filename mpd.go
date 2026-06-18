package main

import (
	"os"
	"fmt"
	"sort"
	"strings"

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

// Get a sorted list of all play lists
func getlists() []string {
	connect_to_mpd()
	lists, err := conn.ListPlaylists()
	if err != nil {
		die(err)
	}

	// Build an array of lists so we can sort them
	l := []string{}
	for _, v := range lists {
		l = append(l, v["playlist"])
	}
	sort.Strings(l)

	return l
}

// Find a play list; if it's not an exact match then try lower case
func find_playlist(list string) string {
	l := getlists()

	// Exact match?
	for _, v := range l {
		if list == v {
			return v
		}
	}

	// Case independent
	list = strings.ToLower(list)
	for _, v := range l {
		if list == strings.ToLower(v) {
			return v
		}
	}

	return ""
}

func load_playlist(list string) {
	found := find_playlist(list)
	if found == "" {
		fmt.Printf("List \"" + list + "\" not found\n")
		return
	}

	conn.Clear()
	err := conn.PlaylistLoad(found, -1, -1)
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
