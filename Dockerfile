FROM gcr.io/moonrhythm-containers/go-scratch

COPY acourse /
COPY template /template
COPY settings /settings
COPY assets /assets
COPY .build/* ./
EXPOSE 8080

ENTRYPOINT ["/acourse"]
