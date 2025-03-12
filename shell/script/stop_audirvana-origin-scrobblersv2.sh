#!/bin/zsh

launchctl stop  com.vincent.audirvana-scrobblerv2.job
launchctl remove  com.vincent.audirvana-scrobblerv2.job
rm -f ~/Library/LaunchAgents/com.vincent.audirvana-origin-scrobblerv2.job.plist
rm -f ./shell/bin/audirvana-origin-scrobbler
rm -f ./shell/bin/nowplaying-cli-mac
sudo rm -f /opt/local/bin/nowplaying-cli-mac



