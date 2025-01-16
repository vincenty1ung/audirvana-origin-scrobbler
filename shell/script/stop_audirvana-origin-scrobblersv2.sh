#!/bin/zsh

launchctl stop  com.vincent.audirvana-scrobblerv2.job
launchctl remove  com.vincent.audirvana-scrobblerv2.job
rm -f ~/Library/LaunchAgents/com.vincent.audirvana-origin-scrobblerv2.job.plist



