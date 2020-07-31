mkdir -p /home/go
mkdir -p /home/go-tools
chown -R $USER:$USER /home/go
chown -R $USER:$USER /home/go-tools

GO_VERSION=1.14.4
GOOS=linux
GOARCH=amd64
GOROOT=/home/go
GOPATH=/home/go-tools

curl -fsSL https://storage.googleapis.com/golang/go$GO_VERSION.$GOOS-$GOARCH.tar.gz | tar -C /home -xzv