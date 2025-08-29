# 🗳️ vote

## features

* zID verification of users
* perform nominations and voting
* multiple client options + email nomination handling

## setup

### clients

**vote** offers two complete client options: 
- a tui (terminal) interface, and
- a web interface

email nominations are automatically handled with a cloudflare email worker 🪄

#### tui

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

the tui interface is made to be served over ssh, and will do so by default. to run locally without ssh, include the following in your `config.yaml` (see [configuration](tui/README.md)):

```yaml
# config.yaml
tui:
  local: true
```

#### web

TODO: web setup guide

### backend

TODO: backend setup guide

---

made with ❤️ by the dev subcom @ [linux society](https://linuxunsw.org/)
