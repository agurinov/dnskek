{{define "agent"}}<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.appName}}</string>

    <key>ProgramArguments</key>
    <array>
        <string>{{executable}}</string>{{if .debug}}
        <string>--debug</string>{{end}}
        <string>serve</string>
    </array>

    <key>StandardOutPath</key>
    <string>{{user.HomeDir}}/Library/Logs/{{.appName}}/access.log</string>
    <key>StandardErrorPath</key>
    <string>{{user.HomeDir}}/Library/Logs/{{.appName}}/error.log</string>

    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:$PATH</string>
    </dict>

    <key>RunAtLoad</key><true/>
</dict>
</plist>{{end}}
