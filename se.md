Client handler

- Tar imot koblinger og legger i state
- Løkke som lytter til innkommende meldinger (Kommandoer)
- Broadcaster Commands til Command handler

Client command reader

- Leser kommandoer fra enkel klient
- Broadcaser til Command Router

Command Router

- Router commands til deres egne go ruites og sender med conn

Egne go rutines

- Håndterer sitt område for en klient
