Spotify filler
==============

Simple executable to add your local music to playlist or library on spotify.

How to use?
-----------
Define environmen variables
- `SPOTIFY_ID`
- `SPOTIFY_SECRET`

With data from your spotify application https://developer.spotify.com/dashboard/applications

Just drag&drop directory on executable or pass directory path as first argument. Application will search for music files with extensions
- `.3gp`
- `.mp3`
- `.aa`
- `.aac`
- `.aif`
- `.au`
- `.flac`
- `.m4a`
- `.ogg`
- `.oga`
- `.wav`
- `.aax`
- `.mp4`

And ask you about playlist which you want to fill with them (Also favourites library).

How it works?
-------------

It just search for file name in spotify API and matches smiliar tracks.
