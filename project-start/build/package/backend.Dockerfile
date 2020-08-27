FROM golang:alpine AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0  go build ./cmd/backend

FROM scratch
WORKDIR /app
COPY --from=build /app/backend /app/backend
COPY --from=build /app/config/backend.* .
ENTRYPOINT ["/app/backend"]