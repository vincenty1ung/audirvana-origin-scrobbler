# audirvana-origin-scrobbler

## v2
### Shell V1
```shell
 sh shell/script/build_audirvana-origin-scrobblers_launchctl.sh
 sh shell/script/start_audirvana-origin-scrobblers.sh
 sh shell/script/[stop_audirvana-origin-scrobblers.sh](shell/script/stop_audirvana-origin-scrobblers.sh)
```
### Go Build
```shell
# your config
vim config/config.yaml 
go build
./audirvana-origin-scrobbler
tail -f logs/go_audirvana-origin-scrobbler.log
```
### Launch
#### MacOS>15.4处理方法
会出现无法读取到roon的播放信息
原因如下
1. https://www.v2ex.com/t/1122960?p=1#reply15
2. https://github.com/TheBoredTeam/boring.notch/issues/417#issuecomment-2805352589 

处理手段按照如下流程启动程序
* https://github.com/Mx-Iris/MediaRemoteWizard
#### 脚本执行
```shell
sh shell/script/build_audirvana-origin-scrobblers_launchctl.sh
sh shell/script/start_audirvana-origin-scrobblers.sh
sh shell/script/stop_audirvana-origin-scrobblers.sh 
```

## todo
* 扫描audirvana db，加入mysql记录，综合分析。
* 支持audirvana 列表自动同步last.fm

## V1
Scrobbles **Audirvana Origin** playing tracks to **Last.fm**. Uses *zsh*, *AppleScript* and *Python 3* ([scrobbler](https://github.com/hauzer/scrobbler/)).

Loops every 3 seconds (the ````DEFAULT_SLEEP_TIME```` variable). The loop time increases to 20 seconds (the ````LONG_SLEEP_TIME```` variable) if Audirvana has been idle for 5 minutes (the ````AUDIRVANA_IDLE_THRESHOLD```` variable).

Scrobbles to Last.fm when 75 % (the ````THRESHOLD```` variable) of the track has been played.

1. Install [scrobbler](https://github.com/hauzer/scrobbler/) (requires Python 3) with ````pip3 install scrobblerh````
2. Authenticate scrobbler to Last.fm with the ````add-user```` command: ````scrobbler add-user````
3. Download this script.
4. Change the ````LAST_FM_USER```` variable in the script to your Last.fm username.
5. Run the script with ````./audirvana-origin-scrobbler.sh````
6. Play music with Audirvana Origin.

Install as service:

1. Set the script executable

2. Create a job ````audirvana-origin-scrobbler.job.plist```` in /library/LaunchAgents/ with this variables:

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
	<key>KeepAlive</key>
	<true/>
	<key>Label</key>
	<string>audirvana-scrobbler.job</string>
	<key>ProgramArguments</key>
	<array>
		<string>/Users/YOUR_USER/audirvana-origin-scrobbler.sh</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
</dict>
</plist>
````

3. Replace YOUR_USER and the environment variables (Python version instalation, etc) with the details of your system.

4. Alternatively use some app like LaunchControl https://www.soma-zone.com/LaunchControl/ to create the job.
