#!/bin/zsh

cd ./clear/
go build
cd ../
mv ./clear/clear ./clear_bin
sudo cp ./com.vincent.audirvana-origin-scrobbler.job.plist ~/Library/LaunchAgents/
sudo cp ./com.vincent.audirvana-origin-scrobbler-clearlogfile.job.plist ~/Library/LaunchAgents/
