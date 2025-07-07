# Tasklist

## TODO

- [ ] BUG: When sending messages it looks a bit buggy since there are no [you] or so to indicate what you sendt
- [ ] BUG: Typing messages are intercepted if a message is recieved. Log messages one step up
- [ ] Refactor and test
- [ ] Github Actions -> Hosting

Kommer melding fra server

- mutex stopper client fra å skrive
- lagrer klientens tekst
- flytter to hakk opp
- printer meldingen
- flytter to hakk ned
- printer det bruker skrev
- låser opp mutex

Problem

- GetCommand venter på enter, så jeg får ikke tak i dataen som er der før enter en trykket på
- Jeg kan jo hente det de har skrevet fra terminalen, men de vil kunne fortsette å skirve uten at jeg får tatt det inn?
- Eller så kan de ødelegge utseende av de andre meldingene plutselig når vi ikke har hørt enter enda
