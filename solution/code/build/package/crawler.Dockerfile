FROM golang:stretch AS build

WORKDIR /code
COPY . .
RUN go build ./cmd/crawler

FROM scratch

WORKDIR /app
COPY --from=build /code/cmd/crawler .

CMD /app/crawler
