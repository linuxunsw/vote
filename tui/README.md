# terminal interface (tui)

## setup

build from source:

```sh
git clone https://github.com/linuxunsw/vote.git
cd ./vote/tui
go build -o vote cmd/vote/main.go
./vote
```

additionally, for development:

```sh
git clone https://github.com/linuxunsw/vote.git
cd ./vote/tui
go run cmd/vote/main.go
```

the tui interface is made to be served over ssh, and will do so by default. to run locally without ssh, include the following in your `config.yaml`:

```yaml
# config.yaml
tui:
  local: true
```

## configuration

configuration of the tui is done via a `config.yaml` file in any of the following directories:

- `$HOME/.config/vote/`
- `$XDG_CONFIG_HOME/vote/`

`config.example.yaml` provides an example configuration.
