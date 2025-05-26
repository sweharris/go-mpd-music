package main

import (
	"fmt"
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

	for _,song := range list {
		fmt.Println(song["file"])
	}
}

var _ bool = register_fn("dblist", dblist, "List all files in MPD database")

func show_dbinfo() {
	list := mpd_db_info()

	for _,song := range list {
		t := song["Time"]
		if t != "" {
			t=" -- " + secs(t)
		}

		fmt.Printf("%s -- %s -- %s -- %s%s\n",song["file"],song["Album"],song["Artist"],song["Title"],t)
	}
}

var _ bool = register_fn("dbinfo", show_dbinfo, "Song details in MPD database")
