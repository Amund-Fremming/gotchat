# Tasklist

Tanker

ClientState

- View
- ClientName
- ServerName

Lobby view

- Connect
- Create
- Exit
- Help

Chat view

- /leave
- if not "/" -> its a message

---

Do better?

- Go rutines inside Run (Room) ??
- State needs mutex lock
- Make CommandDispatcher spawn new routines to handle different actions
- Disconnect handling for both server and client
  - Client may need to request a disconnect, or handle on closing of the socket

---

# At last

- [ ] Refactor and test
- [ ] Github Actions -> Hosting
