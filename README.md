# acourse

[![Build Status](https://travis-ci.org/acoshift/acourse.svg?branch=master)](https://travis-ci.org/acoshift/acourse)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/acourse/badge.svg?branch=master)](https://coveralls.io/github/acoshift/acourse?branch=master)

Acoshift's course system [acourse.io](https://acourse.io)

## Development

### Config

- add config/sql_url to your postgres database
- add config/service_account to your gcloud service account (using for upload storage)
- add gcloud project id to config/project_id
- add gcloud bucket name to config/bucket
- add smtp email config to config/email_{from,password,port,server,user}

### Software

- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [goreload](https://github.com/acoshift/goreload)
- [Node.js LTS](https://nodejs.org/)

### Running

- go get -u github.com/acoshift/goreload
- yarn install
- make style
- make dev

## Testing

```sh
go get -u github.com/stretchr/testify/mock
go get -u github.com/onsi/ginkgo/ginkgo
go get -u github.com/onsi/gomega/...
go test ./...
```

## License

Apache-2.0
