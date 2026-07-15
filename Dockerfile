FROM golang:1.25.6-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /tools-host ./cmd/tools-host

FROM alpine:3.22
RUN adduser -D -H -u 10001 tools
WORKDIR /app
COPY --from=build /tools-host /usr/local/bin/tools-host
COPY --chown=tools:tools public ./public
USER tools
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/tools-host"]
