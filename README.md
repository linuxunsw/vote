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

the tui interface is made to be served over ssh, and will do so by default. 

install with go:

```sh
go install github.com/linuxunsw/vote/tui@latest
```

or, build from source:

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

to run locally without ssh, use flag `--local`

to customise the form titles, set the environment variables `SOCIETY_NAME` and `EVENT_NAME`. if using a `.env` file, check out [the example](https://github.com/linuxunsw/vote/blob/main/tui/.env.example).

#### web

TODO: web setup guide

### backend

TODO: backend setup guide

---

made with ❤️ by the dev subcom @ [linux society](https://linuxunsw.org/)
