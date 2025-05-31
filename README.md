## go-mpd-music

TLDR: This is similar to "mpc" but written for my personal needs

So in the past I used to use iTunes to play my music on my home theater;
I had a Mac Mini plugged into the receiver.  The nice thing about iTunes
is that it has Applescript integration, so I wrote a couple of shell
scripts that would generate applescript to do things; the
[play](https://www.sweharris.org/Scripts/play.txt)
script would let me queue and play music from the command line, and the
[itunes](https://www.sweharris.org/Scripts/itunes.txt) script would let
me manipulate the state of iTunes.

As time went on some functionality no longer worked (eg previously playlists
had shuffle state that could be set; later versions just had a global state
that didn't have an AppleScript control so you needed to emulate clicking
the menu).  This got annoying.  And my Mac Mini is no longer supported, so
the OS is unpatched, and it can't play back ripped BluRays.  It's time
to move on.

Moving to [mpd](https://www.musicpd.org/) on Linux would let me play my
music, but I needed a command line equivalent to my old shell scripts.
`mpc` worked, but wasn't friendly.  I could have written a shell script
around it, but I decided to write a GoLang program instead.

This lets me control `mpd` similar to how I controlled iTunes.

```
% music help
Command options:
        alts: Show alternative options
      dbinfo: Song details in MPD database
      dblist: List all files in MPD database
          dj: Random play of all songs
        goto: Goto mm:ss in current song
        help: This help message
        info: Shows player info
        next: Next track
       pause: Pause
        play: Play; use `-append` to add to the current queue
    playlist: Load a playlist or show current playlist
    previous: Previous track
      random: Random on/off
    realstop: Really stop, not just pause
      repeat: Repeat on/off
      replay: Skip back 5s
      rescan: Full rescan of music directory
      search: Search for a song
     shuffle: Shuffle current playlist
        skip: Skip forward 5s
      status: Shows status of current song (default action)
        stop: Pause (see realstop)
      update: Update index from music directory
```

There are also "alternative" comamnds which are essentially just shorter
options; they're listed seperately to keep the help list short(er).  eg
`music n` would be the same as `music next`

```
% music alts
Alternative shorter commands:
  list => playlist: Load a playlist or show current playlist
     n => next: Next track
     p => previous: Previous track
  prev => previous: Previous track
     s => status: Shows status of current song (default action)
  skip => next: Next track
```

Now I may not be using `mpd` "correctly".  In particular with `dj` mode
I don't use the "random" option; instead I add all the files to the
current playlist and then shuffle that list.  This lets me see the
queue, what has been recently played, and what is upcoming

```
% music list
 -5: The Kinks -- You Really Got Me -- The Greatest No. 1 Singles Disc 1
 -4: Laura Branigan -- Gloria -- Electric Eighties Disc 5
 -3: Abba -- Gonna Sing You My Love Song -- Abba The Originals Disc 2 - The...
 -2: U2 -- Pride (In The Name Of Love) -- the best Rock Album in the world....
 -1: Diana Ross And Lionel Richie -- Endless Love -- The Sound Of Magic Dis...
==== Queen -- Save Me -- The Ultimate Eighties Ballads Songs From The Heart...
  1: Pet Shop Boys -- Yesterday, When I Was Mad -- Smash The Singles 1993-2...
  2: Fun Boy Three -- Tunnel Of Love -- Best Sellers Of The 80's
  3: The Selector -- On My Radio -- Wow That Was The 70's Disc 7
  4: Jean Michel Jarre and Siriusmo -- Circus -- Electronica 2: The Heart O...
  5: Lighthouse Family -- Ocean Drive -- Top Gear - On The Road Again Disc ...
  6: Shakin' Stevens -- Oh Julie -- 25 Years Of Rock And Roll - 1982
  7: Sylvester -- You Make Me Feel (Mighty Real) -- What A Feeling Disc 2
  8: Spandau Ballet -- Gold -- 54 Hits Of The 80's Disc 3
  9: Dusty Springfield -- Nothing Has Been Proved (Dance Mix) -- Now That's...
 10: Paul McCartney -- Say Say Say -- All The Best!
```

This program is not as full featured as `mpc` but it makes day to day
playback simpler (for me, anyway!)

```
% music status
Current status: pause

     Album: The Ultimate Eighties Ballads Songs From The Heart Of The Decade Disc 1
     Track: 1
    Artist: Queen
     Title: Save Me
  Position: 0:05 of 3:48
      File: cvar_theultimateeightiesballads1/01.Save_Me.mp3
```

Many of the commands present the status afterwards

```
% music next
Current status: play

     Album: Smash The Singles 1993-2004
     Track: 5
    Artist: Pet Shop Boys
     Title: Yesterday, When I Was Mad
  Position: 0:00 of 4:02
      File: apetshop_smash2/05.Yesterday,_When_I_Was_Mad.mp3
```

## Play

The `play` command can be used to specify a directory/file in the
library _or_ a file on the filesystem (`mpd` magic handles this).

So I could do `music play /mp3/SONGS/ALBUMS/Annie_Lennox/Medusa/*`
and it will play those songs:

```
% music play /mp3/SONGS/ALBUM/Annie_Lennox/Medusa/*
Current status: play

     Album: Medusa
     Track: 1
    Artist: Annie Lennox
     Title: No More "I Love You's"
  Position: 0:00 of 4:54
      File: /mp3/SONGS/ALBUM/Annie_Lennox/Medusa/01.No_More_I_Love_Yous.mp3
```

Or I can use the entry in my library:

```
% music play alennox_medusa
Current status: play

     Album: Medusa
     Track: 1
    Artist: Annie Lennox
     Title: No More "I Love You's"
  Position: 0:00 of 4:54
      File: alennox_medusa/01.No_More_I_Love_Yous.mp3
```

(the `alennox_medusa` is just how I organise my music; it's what shows in
my music directory).

By default this will replace the existing queue and start playing it
from the start.

The optional parameter `-append` will add these songs to the end of the
queue; `music play -append /mp3/SONGS/ALBUM/Annie_Lennox/Medusa/*`

The optional parameter `-now` will add these songs to the playlist
after the currently playing song, and then will start playing from the
new song.  So, for example if the queue was

```
 -1: Levellers -- Exodus -- TOTP: The Cutting Edge Disc 2
==== Sting & The Police -- Russians -- The Very Best Of Sting
  1: Various -- The Empire Strikes Back (from: Star Wars II) -- Space Theme...
  2: Fat Larry's Band -- Zoom -- Now Yearbook '82 Disc 1
```

and then we did

```
% music play -now /mp3/SONGS/SINGLE/Billy_Joel/We_Didnt_Start_The_Fire/*
```

now the queue would look like

```
 -2: Levellers -- Exodus -- TOTP: The Cutting Edge Disc 2
 -1: Sting & The Police -- Russians -- The Very Best Of Sting
==== Billy Joel -- We Didn't Start The Fire -- We Didn't Start The Fire (CD...
  1: Billy Joel -- Scenes From An Italian Restaurant -- We Didn't Start The...
  2: Billy Joel -- Zanzibar -- We Didn't Start The Fire (CD Single)
  3: Various -- The Empire Strikes Back (from: Star Wars II) -- Space Theme...
  4: Fat Larry's Band -- Zoom -- Now Yearbook '82 Disc 1
...
```

This lets you play a song immediately but then continue with the
existing playlist afterwards.

## Pause

The "stop" command really just does a pause.  If you do a "realstop"
then `mpd` loses track of the play position, so "play" will start from the
beginning of the song.  For this reason I make the "stop" command just do
a "pause" 'cos that's what I normally really want.

So, yeah, this is really just a command line tool optimised for me!

## Connectivity

If `MPD_HOST` is set then it will try to connect to that server
(on `MPD_PORT` or 6600), otherwise it will try the local socket
`/var/run/mpd/socket`

Communication is handled by the `github.com/fhs/gompd/v2` library.
