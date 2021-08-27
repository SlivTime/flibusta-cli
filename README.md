# Go Flibusta CLI

FOR EDUCATIONAL PURPOSES ONLY

Sometimes I find myself in a situation where a book is not presented.
in any online store. There is a chance that what I am looking 
for can be found on Flibusta site. So I made this utility to 
be able to search conviniently.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/slivtime/flibusta-cli)](https://pkg.go.dev/github.com/slivtime/flibusta-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/slivtime/flibusta-cli)](https://goreportcard.com/report/github.com/slivtime/flibusta-cli)
[![Build Status](https://app.travis-ci.com/SlivTime/flibusta-cli.svg?branch=main)](https://app.travis-ci.com/SlivTime/flibusta-cli)
[![codecov](https://codecov.io/gh/SlivTime/flibusta-cli/branch/main/graph/badge.svg?token=OPQGUACUJ5)](https://codecov.io/gh/SlivTime/flibusta-cli)


## Setup
Flibusta is available via Tor-network with Onion routing so we can use [Torproxy](https://github.com/dperson/torproxy)
to grant access. By default it binds http proxy to port 8118. 

All configuration can be done via Environment, but it should work with Torproxy with default ports.

```
> go get -u github.com/slivtime/flibusta-cli@latest
> cd flibusta-cli

# Copy example configuration to where you store your environment. I use .zshenv for it.
> cd $GOPATH/pkg/mod/github.com/slivtime/flibusta-cli@[...version...]  
> cat example.env >> ~/.zshenv
> source ~/.zshenv

# Check
> flibusta-cli search Война и мир
> flibusta-cli info 175105
> flibusta-cli get 175105
```

