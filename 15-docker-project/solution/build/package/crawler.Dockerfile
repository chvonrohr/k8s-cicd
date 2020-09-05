FROM golang:alpine AS build
RUN apk update && apk add ca-certificates tzdata && update-ca-certificates
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 go build -ldflags="-w"  ./cmd/crawler

FROM scratch
WORKDIR /app
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/crawler /app/crawler
COPY --from=build /app/config/crawler.* .
ENTRYPOINT ["/app/crawler"]