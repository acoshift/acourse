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
- Gin
- gRPC
- Protocol Buffers
- Firebase Authentication
- Firebase Storage
- Cloud Datastore
- Docker

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

#### Start

- Use local backend `http://localhost:8080/api`

`make dev`

- Use production backend `https://acourse.io/api`

`make prod`

#### Build

- Build for production

`make build`

- Build for local backend

`make local`

### Backend

#### Required

- Go 1.7.x `brew install go`

```sh
go get github.com/acoshift/acourse/cmd/acourse
cd $GOPATH/src/github.com/acoshift/acourse
```

#### Protoc

`make proto`

#### Start

`make dev`

#### Build Docker

`make docker`

#### Deploy

`make deploy`
