FROM scratch
ENV GOROOT /usr/local
ADD https://golang.org/lib/time/zoneinfo.zip /usr/local/lib/time/
ADD cacert.pem /etc/ssl/certs/ca-certificates.crt
ADD acourse /
COPY templates /templates
COPY private /private
COPY public /public
EXPOSE 8080
CMD ["/acourse"]
