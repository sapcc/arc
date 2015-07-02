set RootDir to POSIX path of ((path to me as text) & "::")
tell application "iTerm"
	activate
	set myterm to (make new terminal)
	tell myterm
		launch session "Default"
		tell the last session
			write text "cd " & RootDir
			write text "mosquitto -c mosquitto.conf"
		end tell
		launch session "Default"
		tell the last session
			write text "cd " & RootDir
			write text "postgres  -D /usr/local/var/postgres"
		end tell
		launch session "Default"
		tell the last session
			write text "cd " & RootDir & "/.."
			write text "sleep 2"
			write text "bin/api-server -e tcp://localhost:1883 -c api-server/db/dbconf.yml"
		end tell
		launch session "Default"
		tell the last session
			write text "cd " & RootDir & "/.."
			write text "scripts/omnitruck.rb"
		end tell
	end tell
end tell