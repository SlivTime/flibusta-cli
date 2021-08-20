# Go Flibusta CLI

FOR EDUCATIONAL PURPOSES ONLY

Sometimes I find myself in a situation where a book is not presented.
in any online store. There is a chance that what I am looking 
for can be found on Flibusta site. So I made this utility to 
be able to search conviniently.

## Setup
Flibusta is available via Tor-network with Onion routing so we can use [Torproxy](https://github.com/dperson/torproxy)
to grant access. By default it binds http proxy to port 8118. 

All configuration can be done via Environment, but it should work with Torproxy with default ports.

```
> git clone git@github.com:SlivTime/flibusta-cli.git
> cd flibusta-cli

# Copy example configuration to where you store your environment. I use .zshenv for it.  
> cat example.env >> ~/.zshenv
> source ~/.zshenv

> go build && go install

# Check
> flibusta search Война и мир
> flibusta get 175105
```

