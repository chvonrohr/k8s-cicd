FROM golang:alpine AS build
WORKDIR /app
RUN apk add gcc g++
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 go build -ldflags="-w" ./cmd/backend

FROM scratch
WORKDIR /app
COPY --from=build /app/backend /app/backend
COPY --from=build /app/config/backend.* .
ENTRYPOINT ["/app/backend"]