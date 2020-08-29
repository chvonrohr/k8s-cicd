FROM golang:alpine AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build ./cmd/crawler
FROM scratch
WORKDIR /app 
COPY --from=build /app/crawler /app/crawler
COPY --from=build /app/config/backend.* .
ENTRYPOINT ["/app/crawler"]