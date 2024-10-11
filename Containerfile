# Builder for loggerd.
FROM docker.io/golang AS builder

COPY . /src
WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/loggerd cmd/loggerd/main.go

# Container for the loggerd web server.
FROM gcr.io/distroless/static AS loggerd

COPY --from=builder /bin/loggerd /bin/loggerd

EXPOSE 8080

VOLUME /log

ENTRYPOINT ["/bin/loggerd"]
