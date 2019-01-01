# vcn-cli

Global, de-centralized certification of digital assets

## Development



Clone this directory into your `GOPATH`, usually this is `$HOME/go/src/`

### PATH

```
$> export GOPATH=$HOME/go
$> export GOBIN=$GOPATH/bin
$> PATH=$PATH:$GOPATH:$GOBIN
$> export PATH
```

or simply put it to `$HOME/.bash_profile` once.

### Installation

Install vcn system-wide.

```
$> cd vcn-cli/vcn
$> go install
```

## Usage

```
$> vcn commit <file> [user]:[password]
$> vcn verify <file>
```

## Distribution

Building the application for different platforms
```
$> env GOOS=windows GOARCH=amd64 go build vcn.go 
```
c.f. [Digital Ocean multi-platform guide](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04)
