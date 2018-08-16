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

## Deployment

After GCB finished (and live in [staging](https://staging.acourse.io)), run `git push origin HEAD:production` to trigger production deployment.

## License

Apache-2.0
