# Setup commands for {{.appName}} on Mac OS X
# We think this is the best option
# But you can act on your own and setup {{.appName}} manually


# Install info
# executable path: {{executable}}
# user:            {{user.Username}} (UID: {{user.Uid}}, GID: {{user.Gid}})


# TODO if not lookup from $PATH -> Phase 0 -> copy to path and chmod


# Phase 1. Prepare OS X launchctl agent

cat > {{user.HomeDir}}/Library/LaunchAgents/{{.appName}}.plist <<AGENT
{{template "agent" .}}
AGENT
mkdir -p {{user.HomeDir}}/Library/Logs/{{.appName}}
launchctl load {{user.HomeDir}}/Library/LaunchAgents/{{.appName}}.plist
# TODO check service is running


# Phase 2. Prepare system DNS resolver

# sudo cat > /etc/resolver/{{.zone}} <<RESOLVER
{{template "resolver" .}}
# RESOLVER
# sudo killall -HUP mDNSResponder


# Run the commands above to install {{.appName}}:
# TODO pass original args (--debug, --port)
# eval "$({{executable}} install)"
