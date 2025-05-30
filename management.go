package main

import (
	"fmt"
	"strings"
)

func rescan() {
	connect_to_mpd()
	job, err := conn.Rescan("")
	if err != nil {
		die(err)
	}
	fmt.Println("Rescan started; job ID", job)
}

var _ bool = register_fn("rescan", rescan, "Full rescan of music directory")

func update() {
	connect_to_mpd()
	job, err := conn.Update("")
	if err != nil {
		die(err)
	}
	fmt.Println("Update started; job ID", job)
}

var _ bool = register_fn("update", update, "Update index from music directory")

func dblist() {
	list := mpd_db_info()

	for _, song := range list {
		fmt.Println(song["file"])
	}
}

var _ bool = register_fn("dblist", dblist, "List all files in MPD database")

func show_dbinfo() {
	list := mpd_db_info()

	for _, song := range list {
		t := song["Time"]
		if t != "" {
			t = " -- " + secs(t)
		}

		fmt.Printf("%s -- %s -- %s -- %s%s\n", song["file"], song["Album"], song["Artist"], song["Title"], t)
	}
}

var _ bool = register_fn("dbinfo", show_dbinfo, "Song details in MPD database")

// Searches artist/title/album and returns the file for matching entries
func search_db() {
	if len(Args) == 0 {
		die("Missing search string")
	}
	// srch_str := strings.ToLower(Args[0])
	srch_str := strings.ToLower(strings.Join(Args, " "))

	list := mpd_db_info()

	for _, song := range list {
		if strings.Contains(strings.ToLower(song["Artist"]), srch_str) ||
			strings.Contains(strings.ToLower(song["Album"]), srch_str) ||
			strings.Contains(strings.ToLower(song["Title"]), srch_str) {
			fmt.Println(song["file"])
		}
	}
}

var _ bool = register_fn("search", search_db, "Search for a song")
