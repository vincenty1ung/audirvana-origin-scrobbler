#!/bin/zsh

go build
mv ./lastfm-scrobbler ./shell/bin/lastfm-scrobbler
cd ./roon/now-playing/now-playing || exit
make
cd ../../../
sudo cp ./roon/now-playing/now-playing/nowplaying-cli-mac /opt/local/bin/nowplaying-cli-mac
# mv ./roon/now-playing/now-playing/nowplaying-cli-mac "$GOBIN"/nowplaying-cli-mac
mv ./roon/now-playing/now-playing/nowplaying-cli-mac ./shell/bin/nowplaying-cli-mac
sudo cp ./shell/launch/com.vincent.lastfm-scrobbler.job.plist ~/Library/LaunchAgents/
