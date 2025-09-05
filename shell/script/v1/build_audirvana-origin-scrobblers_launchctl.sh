#!/bin/zsh

cd ./shell/clear
go build
cd ../../
mv ./shell/clear/clear ./shell/bin/clear_bin
sudo cp ./shell/launch/com.vincent.lastfm-scrobbler.job.plist ~/Library/LaunchAgents/
sudo cp ./shell/launch/com.vincent.lastfm-scrobbler-clearlogfile.job.plist ~/Library/LaunchAgents/
