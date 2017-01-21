# Acourse

[![Standard - JavaScript Style Guide](https://img.shields.io/badge/code%20style-standard-brightgreen.svg)](http://standardjs.com/)

Acoshift's course system [Link](https://acourse.io)

## Stack

### Frontend

- Vue.js
- Rxjs + VueRx
- Webpack

### Backend

- Go
- gRPC
- Protocol Buffers
- Firebase Authentication
- Firebase Storage
- Cloud Datastore
- Docker

---

## Development

### Frontend

#### Required

- Node.js
- yarn `npm install -g yarn`

```sh
git clone https://github.com/acoshift/acourse.git
cd acourse/ui
yarn
```

#### Development

- Run backend on localhost (load config from private/config.stag.yaml)

  - Start backend by run `make dev`

  - Start ui on port 9000 by run `cd ui && make dev`

- Use production backend

  - Start ui on port 9000 by run `cd ui && make prod`

#### Build

- For build in production mode, see on backend section

- For build for local backend (for test preload data)

  - Run `make local` on root directory

### Backend

#### Required

- Go 1.7.x `brew install go`

```sh
go get github.com/acoshift/acourse/cmd/acourse
cd $GOPATH/src/github.com/acoshift/acourse
```

#### Generate Protocol Buffers

- Install protoc and protoc-gen-go

- Run `make proto`

#### Development

`make dev`

#### Deploy

`make deploy` or just `make`
