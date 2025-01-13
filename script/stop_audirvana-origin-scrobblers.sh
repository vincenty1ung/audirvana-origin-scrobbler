#!/bin/zsh

launchctl stop  com.vincent.audirvana-scrobbler.job
launchctl remove  com.vincent.audirvana-scrobbler.job
rm -f ~/Library/LaunchAgents/com.vincent.audirvana-origin-scrobbler.job.plist

launchctl stop  com.vincent.audirvana-scrobbler-clearlogfile.job
launchctl remove  com.vincent.audirvana-scrobbler-clearlogfile.job
rm -f ~/Library/LaunchAgents/com.vincent.audirvana-origin-scrobbler-clearlogfile.job.plist


