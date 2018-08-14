FROM acoshift/go-scratch

ADD acourse /
COPY template /template
COPY settings /settings
COPY static /static
COPY static.yaml /static.yaml
EXPOSE 8080

ENTRYPOINT ["/acourse"]
