FROM scratch
MAINTAINER Thanatat Tamtan <acoshift@gmail.com>

# Setup Go
ENV GOROOT /usr/local
ADD https://golang.org/lib/time/zoneinfo.zip /usr/local/lib/time/
ADD cacert.pem /etc/ssl/certs/ca-certificates.crt

# Setup App
ADD acourse /
COPY config.yaml /
COPY acourse_io.crt /
COPY acourse_io.key /
EXPOSE 80
EXPOSE 443
ENTRYPOINT ["/acourse"]
