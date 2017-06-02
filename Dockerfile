FROM acoshift/go-scratch

ADD entrypoint /
COPY template /template
COPY static /static
EXPOSE 8080

ENTRYPOINT ["/entrypoint"]
