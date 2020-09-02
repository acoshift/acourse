你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# acourse

[![Build Status](https://travis-ci.org/acoshift/acourse.svg?branch=master)](https://travis-ci.org/acoshift/acourse)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/acourse/badge.svg?branch=master)](https://coveralls.io/github/acoshift/acourse?branch=master)

Acourse Website [acourse.io](https://acourse.io)

## Development

### Config

- add config/sql_url to your cockroachdb database
- add config/service_account to your gcloud service account (using for upload storage)
- add gcloud project id to config/project_id
- add gcloud bucket name to config/bucket
- add smtp email config to config/email_{from,password,port,server,user}

### Software

- [CockroachDB](https://www.cockroachlabs.com/)
- [Redis](https://redis.io/) (optional)
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
