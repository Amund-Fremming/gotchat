module client

require github.com/amund-fremming/common v0.0.0

require github.com/gorilla/websocket v1.5.3

require github.com/chzyer/readline v1.5.1

require (
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace github.com/amund-fremming/common => ../common

go 1.24.4
