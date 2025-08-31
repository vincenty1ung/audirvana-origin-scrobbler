#!/bin/zsh

launchctl stop  com.vincent.audirvana-scrobbler.job
launchctl remove  com.vincent.audirvana-scrobbler.job
rm -f ~/Library/LaunchAgents/com.vincent.audirvana-origin-scrobbler.job.plist
rm -f ./shell/bin/audirvana-origin-scrobbler
rm -f ./shell/bin/nowplaying-cli-mac
sudo rm -f /opt/local/bin/nowplaying-cli-mac



