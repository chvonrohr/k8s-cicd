mkdir -p /home/go
mkdir -p /home/go-tools
chown -R $USER:$USER /home/go
chown -R $USER:$USER /home/go-tools

curl https://sdk.cloud.google.com > install.sh
bash install.sh --disable-prompts
rm install.sh
gcloud components install kubectl kustomize alpha beta --quiet

apt-get update
apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"
apt-get update
apt-get install docker-ce-cli
