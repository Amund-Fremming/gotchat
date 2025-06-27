# Tasklist

- Go rutines inside Run (Room) ??
- State needs mutex lock
- Make CommandDispatcher spawn new routines to handle different actions
- Disconnect handling for both server and client
  - Client may need to request a disconnect, or handle on closing of the socket

---

## Bugs

- [ ] No error handling, removing when a write operation fails
  - Make a error handling function that closes and removes client
- [ ] If you are typing a message, and you recieve a message, you ui gets cronked
  - Probably just log one step up on all messages

## At last

- [ ] Refactor and test
- [ ] Github Actions -> Hosting
