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

clone the repo and create a `.env` file (see [the example](https://github.com/linuxunsw/vote/blob/main/tui/.env.example)). to build:

```bash
cd tui
go build -o vote cmd/vote/main.go
./vote
```

or, install with go:

```bash
go install github.com/linuxunsw/vote/tui@latest
```

for development:

```bash
cd tui
go run cmd/vote/main.go
```

TODO: add hosting information

#### web

TODO: web setup guide

### backend

TODO: backend setup guide

---

made with ❤️ by the dev subcom @ [linux society](https://linuxunsw.org/)
