FROM golang:1.16.0

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN mkdir -p /workspace
WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -trimpath -o acourse -ldflags "-w -s" .

FROM node:14.15.5

ENV NODE_ENV=production

RUN mkdir -p /workspace
WORKDIR /workspace
ADD package.json yarn.lock ./
RUN yarn install
ADD . .
RUN yarn run gulp

FROM gcr.io/moonrhythm-containers/go-scratch

COPY --from=0 /workspace/acourse /
COPY --from=0 /workspace/template /template
COPY --from=0 /workspace/settings /settings
COPY --from=1 /workspace/assets /assets
COPY --from=1 /workspace/.build/* /

EXPOSE 8080

ENTRYPOINT ["/acourse"]
