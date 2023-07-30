 
# Subtitle matcher

This project will rename your subtitle files to match the video file.

# Requirements
go version >=1.20.2

# Build
```
go build
```

# Usage
Example files to match
```
$ ls -l __test
MyMovie_1.srt
MyMovie.avi
MySerial S01E01.avi
MySerial S01E01_x254.srt
```
## Running the subtitle matcher
```
$ ./subMatcher ./__test
--- Movies
MyMovie.avi
MySerial S01E01.avi
--- Subs
0 MyMovie_1.srt
1 MySerial S01E01_x254.srt
--- --- ---
score  5 MyMovie.avi MyMovie_1.srt
----
Matched subs forMyMovie.avi
Rename
./__test/MyMovie_1.srt
to
./__test/MyMovie.srt
[Y/n] 
score  13 MySerial S01E01.avi MySerial S01E01_x254.srt
----
Matched subs forMySerial S01E01.avi
Rename
./__test/MySerial S01E01_x254.srt
to
./__test/MySerial S01E01.srt
[Y/n]
```

## Result
```
$ ls ./__test
MyMovie.avi  MyMovie.srt  MySerial S01E01.avi  MySerial S01E01.srt
```
# Options
## `-r`
Recursive option. Every directory will have `subMatcher` applied for itself.
