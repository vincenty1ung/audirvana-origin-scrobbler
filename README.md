# audirvana-origin-scrobbler

Scrobbles **Audirvana Origin** playing tracks to **Last.fm**. Uses *zsh*, *AppleScript* and *Python 3* ([scrobbler](https://github.com/hauzer/scrobbler/)).

Loops every 3 seconds (the ````DEFAULT_SLEEP_TIME```` variable). The loop time increases to 20 seconds (the ````LONG_SLEEP_TIME```` variable) if Audirvana has been idle for 5 minutes (the ````AUDIRVANA_IDLE_THRESHOLD```` variable).

Scrobbles to Last.fm when 75 % (the ````THRESHOLD```` variable) of the track has been played.

1. Install [scrobbler](https://github.com/hauzer/scrobbler/) (requires Python 3) with ````pip3 install scrobblerh````
2. Authenticate scrobbler to Last.fm with the ````add-user```` command: ````scrobbler add-user````
3. Download this script.
4. Change the ````LAST_FM_USER```` variable in the script to your Last.fm username.
5. Run the script with ````sh ./audirvana-origin-scrobbler.sh````
6. Play music with Audirvana Origin.

Install as service:

1. Set the script executable

2. Create a job in /library/LaunchAgents/ with this variables:

````
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>EnvironmentVariables</key>
	<dict>
		<key>PATH</key>
		<string>/usr/local/sbin:/opt/local/bin:/opt/local/sbin:/Library/Frameworks/Python.framework/Versions/3.10/bin:/Library/Frameworks/Python.framework/Versions/3.9/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/Users/YOUR_USER/Library/Python/3.10/bin:/Library/Apple/usr/bin</string>
	</dict>
	<key>Label</key>
	<string>audirvana-scrobbler.job</string>
	<key>ProgramArguments</key>
	<array>
		<string>sh</string>
		<string>/Users/YOUR_USER/audirvana-origin-scrobbler.sh</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
</dict>
</plist>
````

3. Replace in the file YOUR_USER and the environment variables (Python version instalation, etc) with the details of your system.

4. Alternatively use some app like LaunchControl https://www.soma-zone.com/LaunchControl/ to create the job.
