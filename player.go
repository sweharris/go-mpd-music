package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func show_status() {
	status := get_status()
	song := get_song()

	s := status["state"]

	fmt.Printf("Current status: %s\n", s)
	fmt.Println()
	fmt.Printf("     Album: %s\n", song["Album"])
	fmt.Printf("     Track: %s\n", song["Track"])
	fmt.Printf("    Artist: %s\n", song["Artist"])
	fmt.Printf("     Title: %s\n", song["Title"])
	if s != "stop" {
		fmt.Printf("  Position: %s of %s\n", secs(status["elapsed"]), secs(status["duration"]))
	}
	fmt.Printf("      File: %s\n", song["file"])
}

var _ bool = register_fn("status", show_status, "Shows status of current song (default action)")
var _ bool = register_alt("s", "status")

func json_status() {
	type Status struct {
		State    string `json:"state"`
		Album    string `json:"album"`
		Track    string `json:"track"`
		Artist   string `json:"artist"`
		Title    string `json:"title"`
		Elapsed  string `json:"elapsed"`
		Duration string `json:"duration"`
		File     string `json:"file"`
	}

	status := get_status()
	song := get_song()
	s := status["state"]

	j := Status{
		State:    s,
		Album:    song["Album"],
		Track:    song["Track"],
		Artist:   song["Artist"],
		Title:    song["Title"],
		Elapsed:  secs(status["elapsed"]),
		Duration: secs(status["duration"]),
		File:     song["file"],
	}

	str, _ := json.Marshal(j)
	fmt.Println(string(str))
}

var _ bool = register_fn("jsonstatus", json_status, "Show status in JSON format")

func show_info() {
	status := get_status()
	song := get_song()
	fmt.Printf("    State = %s\n", status["state"])
	fmt.Printf("   Repeat = %s\n", show_truefalse(status["repeat"]))
	fmt.Printf("   Random = %s\n", show_truefalse(status["random"]))
	fmt.Printf("  Consume = %s\n", show_truefalse(status["consume"]))
	fmt.Printf("   SongID = %s (%s)\n", status["songid"], song["file"])
	if status["updating_db"] != "" {
		fmt.Println("Database is being updated; job ID", status["updating_db"])
	}
}

var _ bool = register_fn("info", show_info, "Shows player info")

func play() {
	status := get_status()

	// If there are no arguments, just tell the daemon to play
	// If we've stopped we need to start the find the current song
	// and start that, otherwise turn off pause
	if len(Args) == 0 {
		if status["state"] == "stop" {
			song := get_song()
			p, _ := strconv.Atoi(song["Pos"])
			conn.Play(p)
		} else {
			conn.Pause(false)
		}
	} else {
		// We can -append to the current list, or replace it
		append := false
		now := false
		// We want to play a selected file
		if Args[0] == "-append" {
			Args = Args[1:]
			append = true
		} else if Args[0] == "-now" {
			Args = Args[1:]
			now = true
		} else {
			conn.Clear()
		}
		// We don't want to play these in random order and we do
		// want to repeat.
		conn.Random(false)
		conn.Repeat(true)
		// We're not doing any sanity checks here, so we can
		// specify filenames or relative to the music directory
		// If we specify filenames we can play songs without adding
		// them to the catalogue.
		for _, n := range Args {
			e := conn.Add(n)
			if e != nil {
				die(e)
			}
		}
		if !append {
			conn.Play(0)
		}

		// If we specified -now then move these new songs to the
		// front of the queue
		if now {
			new_status := get_status()
			old_len, _ := strconv.Atoi(status["playlistlength"])
			new_len, _ := strconv.Atoi(new_status["playlistlength"])
			pos, _ := strconv.Atoi(status["song"])
			if new_len > old_len {
				conn.Move(old_len, new_len, pos+1)
				// Huh, moving a song moves the play position
				// to 0, so Next() doesn't work.  Play from
				// this position instead
				conn.Play(pos + 1)
			}
		}
	}

	show_status()
}

var _ bool = register_fn("play", play, "Play; use `-append` to add to the current queue")

func pause() {
	connect_to_mpd()
	conn.Pause(true)
	show_status()
}

var _ bool = register_fn("pause", pause, "Pause")

// I don't really like "stop" so make it a pause
var _ bool = register_fn("stop", pause, "Pause (see realstop)")

func stop() {
	connect_to_mpd()
	conn.Stop()
	show_status()
}

var _ bool = register_fn("realstop", stop, "Really stop, not just pause")

func next() {
	connect_to_mpd()
	conn.Next()
	show_status()
}

var _ bool = register_fn("next", next, "Next track")
var _ bool = register_alt("skip", "next")
var _ bool = register_alt("n", "next")

func previous() {
	connect_to_mpd()
	conn.Previous()
	show_status()
}

var _ bool = register_fn("previous", previous, "Previous track")
var _ bool = register_alt("prev", "previous")
var _ bool = register_alt("p", "previous")

func skip_forward() {
	connect_to_mpd()
	conn.SeekCur(time.Duration(5*time.Second), true)
	show_status()
}

var _ bool = register_fn("skip", skip_forward, "Skip forward 5s")

func skip_back() {
	connect_to_mpd()
	conn.SeekCur(time.Duration(-5*time.Second), true)
	show_status()
}

var _ bool = register_fn("replay", skip_back, "Skip back 5s")

func goto_time() {
	if len(Args) == 0 {
		die("Missing time")
	}
	t := Args[0]
	if strings.Contains(t, ":") {
		// Time has been specified as m:ss
		// convert to ##m##s format
		t_a := strings.SplitN(t, ":", 2)
		if t_a[0] == "" {
			t_a[0] = "0"
		}
		if t_a[1] == "" {
			t_a[1] = "0"
		}
		t = t_a[0] + "m" + t_a[1] + "s"
	}
	d, err := time.ParseDuration(t)
	if err != nil {
		die(err)
	}
	connect_to_mpd()
	conn.SeekCur(d, false)
	show_status()
}

var _ bool = register_fn("goto", goto_time, "Goto mm:ss in current song")

// Take the list of songs in the database, make this the playlist,
// shuffle it, and play.
func dj() {
	connect_to_mpd()
	conn.Clear()
	conn.Random(false)
	conn.Repeat(true)
	conn.Add("")
	conn.Shuffle(-1, -1)
	conn.Play(0)
	show_status()
}

var _ bool = register_fn("dj", dj, "Random play of all songs")

// Without a parameter show the current playlist -5/+10 items, otherwise
// load the named playlist and play it
func playlist() {
	connect_to_mpd()
	// If there is an argument and it's a number then use this as
	// the number of entries +/- to display.
	before := 5
	after := 10
	if len(Args) == 1 && Args[0][:1] == "-" {
		to_display, e := strconv.Atoi(Args[0])
		// If the conversion worked correctly...
		if e == nil {
			// to_display will be negative (eg -5) so we
			// make it positive here
			before = -to_display
			after = -to_display

			// We've consumed the only entry
			Args = nil
		}
	}
	if len(Args) == 0 {
		// Display the current playlist; previous 5 songs and the
		// next 10.  Highlight the current song
		song := get_song()
		p := song["Pos"]
		pos, _ := strconv.Atoi(p)
		start := pos - before  // songs before current position
		end := pos + after + 1 // songs after current position
		if start < 0 {
			start = 0 // We don't have 5 earlier songs!
		}
		playlist := get_playlist(start, end)

		// playlist appears to be in order, so this
		// loop is easy
		for _, plsong := range playlist {
			this_pos, _ := strconv.Atoi(plsong["Pos"])

			// Position marker, -5 to 10
			p_str := fmt.Sprintf("%3d:", this_pos-pos)
			if this_pos == pos {
				p_str = "===="
			}

			// Try to make songs without tags a little prettier
			a := plsong["Artist"]
			if a != "" {
				a += " -- "
			}
			t := plsong["Title"]
			if t == "" {
				t = plsong["file"]
			}
			al := plsong["Album"]
			if al != "" {
				al = " -- " + al
			}

			str := fmt.Sprintf("%s %s%s%s", p_str, a, t, al)
			if len(str) > 75 {
				str = str[:75] + "..."
			}
			fmt.Println(str)
		}
	} else {
		// We want to load a new playlist

		// We don't want to play these in random order and we do
		// want to repeat.
		conn.Random(false)
		conn.Repeat(true)
		load_playlist(Args[0])
		// Play from the start of the list
		conn.Play(0)
		show_status()
	}
}

var _ bool = register_fn("playlist", playlist, "Load a playlist or show current playlist")
var _ bool = register_alt("list", "playlist")

func showlist() {
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

	fmt.Println("Available playlists:")

	for _, v := range l {
		// We print this with " " around it for simple cut'n'paste
		fmt.Printf("  \"%s\"\n", v)
	}
}

var _ bool = register_fn("showlist", showlist, "Show available playlists")
var _ bool = register_alt("showlists", "showlist")

func repeat() {
	if len(Args) == 0 {
		die("Need on/off/true/false")
	}
	connect_to_mpd()
	conn.Repeat(truefalse(Args[0]))
	show_info()
}

var _ bool = register_fn("repeat", repeat, "Repeat on/off")

func random() {
	if len(Args) == 0 {
		die("Need on/off/true/false")
	}
	connect_to_mpd()
	conn.Random(truefalse(Args[0]))
	show_info()
}

var _ bool = register_fn("random", random, "Random on/off")

func shuffle() {
	connect_to_mpd()
	conn.Shuffle(-1, -1)
}

var _ bool = register_fn("shuffle", shuffle, "Shuffle current playlist")
