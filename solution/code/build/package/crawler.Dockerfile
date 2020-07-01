# we use stretch linux to build
FROM golang:stretch AS build 

WORKDIR /code
# copy dependency files first
# this way docker can skip the go mod download if go.mod or go.sum didn't change
COPY go.mod .
COPY go.sum .
# install go dependencies
RUN go mod download

# Todo: ssl cert

# copy everything from current working directory into image
COPY . .
RUN go build ./cmd/crawler

FROM scratch

WORKDIR /app
COPY --from=build /code/crawler .

CMD /app/crawler
