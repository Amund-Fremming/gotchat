module client

require github.com/amund-fremming/common v0.0.0

require github.com/gorilla/websocket v1.5.3

require (
	github.com/chzyer/readline v1.5.1 
	github.com/joho/godotenv v1.5.1 
	golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5 
)

replace github.com/amund-fremming/common => ../common

go 1.24.4
