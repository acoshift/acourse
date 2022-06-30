FROM node:14

ENV NODE_ENV=production

RUN mkdir -p /workspace
WORKDIR /workspace
ADD package.json yarn.lock ./
RUN yarn install
ADD . .
RUN yarn run gulp

FROM golang:1.18.3

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN mkdir -p /workspace
WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download
COPY --from=0 /workspace/ ./
RUN go build -trimpath -o .build/acourse -ldflags "-w -s" .

FROM gcr.io/distroless/static

COPY --from=1 /workspace/.build/* /

EXPOSE 8080
ENTRYPOINT ["/acourse"]
