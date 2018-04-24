# Soccer Robot Remote Development

## Development Environment

The development environment heavily relies on [`make`](https://www.gnu.org/software/make/). Under the hood, `make` calls other tools such as `go`, `yarn` etc. Let's first make sure you have `go`, `node` and `yarn`:

On macOS using [Homebrew](https://brew.sh):

```sh
brew install go node yarn
```

On Ubuntu (or Ubuntu [using the Windows Subsystem for Linux](https://www.microsoft.com/nl-NL/store/p/ubuntu/9nblggh4msv6?rtc=1)):

```sh
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list

curl -sS https://deb.nodesource.com/gpgkey/nodesource.gpg.key | sudo apt-key add -
echo "deb https://deb.nodesource.com/node_8.x xenial main" | sudo tee /etc/apt/sources.list.d/nodesource.list
echo "deb-src https://deb.nodesource.com/node_8.x xenial main" | sudo tee -a /etc/apt/sources.list.d/nodesource.list

sudo apt-get update
sudo apt-get install build-essential nodejs yarn

curl -sS https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz -o go1.10.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.10.1.linux-amd64.tar.gz
sudo ln -s /usr/local/go/bin/* /usr/local/bin
```

### Getting started with Go Development

We will first need a Go workspace. The Go workspace is a folder that contains the following sub-folders:

- `src` which contains all source files
- `pkg` which contains compiled package objects
- `bin` which contains binary executables

From now on this folder is referred to as `$GOPATH`. By default, Go assumes that it's in `$HOME/go`, but to be sure that everything works correctly, you can add the following to your profile (in `$HOME/.profile`):

```sh
export GOPATH="$(go env GOPATH)"
export PATH="$PATH:$GOPATH/bin"
```

Now that your Go development environment is ready, it strongly recommended to get familiar with Go by following the [Tour of Go](https://tour.golang.org/).

### Getting started with development
_Note the `--recursive` flag!_
```sh
git clone --recursive git@github.com:rvolosatovs/turtlitto.git
```

All development is done in this directory.

```sh
cd $GOPATH/src/github.com/rvolosatovs/turtlitto
```

As most of the tasks will be managed by `make` we will first initialize the tooling. You might want to run this command from time to time:

```sh
make deps
```

#### Folder Structure

```
.
├── STYLE.md     guidelines for contributing: branching, commits, code style, etc.
├── DEVELOPMENT.md      guide for setting up your development environment
├── Gopkg.lock          dependency lock file managed by golang/dep
├── Gopkg.toml          dependency file managed by golang/dep
├── Makefile            dev/test/build tooling
├── README.md           general information about this project
│   ...
├── cmd                 contains the different binaries
│   └── soccer-robot-remote          contains the Soccer Robot Remote
├── docs                contains the documentation
├── front               contains the frontend of the project
├── pkg                 contains all libraries used in the backend
├── release             binaries will be compiled to this folder - not added to git
└── vendor              dependencies managed by golang/dep - not added to git
```

#### Testing

For backend:
```sh
make go.test
```

For frontend:
```sh
make js.test
```

#### Building

There's one binary to be built: the `soccer-robot-remote` binary, which holds the remote control for the soccer robots.

To build it run:

```
make soccer-robot-remote
```

This will result in `release/soccer-robot-remote-linux-amd64` generated (suffix can differ based on your architecture and operating system).
