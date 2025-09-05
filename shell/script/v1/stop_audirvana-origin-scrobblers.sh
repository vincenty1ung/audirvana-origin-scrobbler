#!/bin/zsh

launchctl stop  com.vincent.lastfm-scrobbler.job
launchctl remove  com.vincent.lastfm-scrobbler.job
rm -f ~/Library/LaunchAgents/com.vincent.lastfm-scrobbler.job.plist

launchctl stop  com.vincent.lastfm-scrobbler-clearlogfile.job
launchctl remove  com.vincent.lastfm-scrobbler-clearlogfile.job
rm -f ~/Library/LaunchAgents/com.vincent.lastfm-scrobbler-clearlogfile.job.plist


