FROM scratch
MAINTAINER Thanatat Tamtan <acoshift@gmail.com>

# Setup Go
ENV GOROOT /usr/local
ADD https://golang.org/lib/time/zoneinfo.zip /usr/local/lib/time/
ADD cacert.pem /etc/ssl/certs/ca-certificates.crt

# Setup App
ADD acourse /
COPY templates /templates
COPY private /private
COPY public /public
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/acourse"]
