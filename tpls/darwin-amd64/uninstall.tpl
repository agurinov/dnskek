# Uninstall commands for dnskek on Mac OS X
# We think this is the best option
# But you can act on your own and uninstall dnskek manually


# Uninstall info
# executable path: {{executable}}
# user:            {{user.Username}} (UID: {{user.Uid}}, GID: {{user.Gid}})


# Phase 1. Rm OS X launchctl agent
launchctl unload {{user.HomeDir}}/Library/LaunchAgents/{{.appName}}.plist
launchctl remove {{.appName}}
rm -f {{user.HomeDir}}/Library/LaunchAgents/{{.appName}}.plist
rm -rf {{user.HomeDir}}/Library/Logs/{{.appName}}


# Run this command to uninstall dnskek:
# TODO pass original args (--debug, --port)
# eval "$({{executable}} uninstall)"
