# Request Logger

Request logger is a web server that logs all requests to stdout.

## Usage

```bash
podman run -p 8080:8080 ghcr.io/tslnc04/loggerd -p 8080
```

Or, if using Docker:

```bash
docker run -p 8080:8080 ghcr.io/tslnc04/loggerd -p 8080
```

To see the help message:

```bash
podman run ghcr.io/tslnc04/loggerd -help
```

## Motivation

I see a lot of random requests in the logs to my services and want to see more about them. In particular, since the services are behind a reverse proxy, they never correctly log the origin of the request from the X-Forwarded-For header. This services allows me to see all the information about requests as they come in and hopefully satisfy my curiosity.

## Copyright

Copyright 2024 Kirsten Laskoski

Licensed under MIT. See [LICENSE](LICENSE) for details.
