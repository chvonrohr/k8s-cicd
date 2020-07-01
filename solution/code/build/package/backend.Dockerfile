# we use stretch linux to build
FROM golang:stretch AS build 

WORKDIR /code
# copy dependency files first
# this way docker can skip the go mod download if go.mod or go.sum didn't change
COPY go.mod .
COPY go.sum .
# install go dependencies
RUN go mod download
# copy everything from current working directory into image
COPY . .
RUN go build ./cmd/backend

# we use "scratch" image to run go service
# the scratch image "doesn't contain anything"
FROM scratch

WORKDIR /app
COPY --from=build /code/backend .

CMD /app/backend