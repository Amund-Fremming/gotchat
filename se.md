Connection/Client handler

- Tar imot koblinger og legger i state
- Løkke som lytter til innkommende meldinger (Kommandoer)
- Broadcaster Commands til Command handler

Command handler

- Tolker commandoer
- Broadcast til
  - Join kanal
  - Leave room kanal
  - Leave program kanal
  - Broadcast melding kanal
  - Info kanal

Join room

- Legger til bruker i rommet
- Sender ny state til klient

Send message

- broadcaster meldingen til rommet

Leave room

- Sender ny state til klient
- Fjerer bruker fra state
- Stopper koblingen

Leave program

- Fjerner bruker fra connections
- Sende ny state
