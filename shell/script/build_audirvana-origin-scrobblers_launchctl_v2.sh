#!/bin/zsh

go build
mv ./audirvana-origin-scrobbler ./shell/bin/audirvana-origin-scrobbler
sudo cp ./shell/launch/com.vincent.audirvana-origin-scrobblerv2.job.plist ~/Library/LaunchAgents/
