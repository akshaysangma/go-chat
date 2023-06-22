# go-chat
# Chat Room written in Golang

# Get Started
## To build the server binary
```sh
cd server
make build
```

## To run the server locally with minimal effort

Ensure you have go toolchain (1.18.0+) and docker running and run the below make command to start a local posgresql database with go-chat database with migration applied.

```sh
# TODO: merge multi-line make command to single make command
cd server
make postgresinit 
make createdb
make migrateup
make build
./bin/server.exe
```

