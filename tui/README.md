# terminal interface (tui)

## setup

build from source:

```sh
git clone https://github.com/linuxunsw/vote.git
cd ./vote/tui
go build -o vote-tui cmd/vote/main.go
./vote-tui
```

additionally, for development:

```sh
git clone https://github.com/linuxunsw/vote.git
cd ./vote/tui
go run cmd/vote/main.go
```

## configuration

configuration of the tui is done via a `config.yaml` file in any of the following directories:

- `$HOME/.config/vote/`
- `$XDG_CONFIG_HOME/vote/`

`config.example.yaml` provides an example configuration. to test the example, set `server` to the backend, and connect using:

```sh
ssh localhost -p 2222
```

## running on privileged ports

to allow users to connect to `vote` without specifying a port, the tui interface needs to be configured to use port 22:

```yaml
# config.yaml
tui:
  port: 22
```

then change the default ssh port:

```
# /etc/ssh/sshd_config
Port 2222 
```

*check that you can connect via this port before continuing!*

to allow the tui to bind to privileged ports (port 22), build and give permissions to bind to privileged ports:

```sh
go build -o vote-tui cmd/vote/main.go
sudo setcap CAP_NET_BIND_SERVICE=+eip vote-tui
./vote-tui

# then, connect from another terminal
ssh user@your.domain.here
```
