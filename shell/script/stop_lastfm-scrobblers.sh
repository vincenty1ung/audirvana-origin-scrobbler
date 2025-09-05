#!/bin/zsh

launchctl stop  com.vincent.lastfm-scrobbler.job
launchctl remove  com.vincent.lastfm-scrobbler.job
rm -f ~/Library/LaunchAgents/com.vincent.lastfm-scrobbler.job.plist
rm -f ./shell/bin/lastfm-scrobbler
rm -f ./shell/bin/nowplaying-cli-mac
sudo rm -f /opt/local/bin/nowplaying-cli-mac



