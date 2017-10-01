# Yarder

a log storage system

# Install

## From source

* Install golang 1.7, 1.8 or 1.9 from https://golang.org/dl/
* Run `go get github.com/log-yarder/yarder`

# Development

* Install golang 1.7, 1.8 or 1.9 from https://golang.org/dl/
* Setup a workspace (see https://golang.org/doc/code.html)
* Run: 
  ```
  mkdir $GOPATH/src/github.com/log-yarder
  git clone https://github.com/log-yarder/yarder $GOPATH/src/github.com/log-yarder/yarder
  cd $GOPATH/src/github.com/log-yarder
  make reset
  make test
  ```
