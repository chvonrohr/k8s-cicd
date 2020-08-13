
# Installation

Thank you for participating in a letsboot.com Kuberentes training.

For everyone:
* gitlab.com user 
* fill in our preparation survey

## Letsboot Virtual Desktop

* We provide you a completly browser based (Theia) virtual desktop environment with all tools, configurations, permissions and the course material. 

* The only thing you need is a working laptop with Chrome and a gitlab user. (On request in advance, we can provide a training laptop.)

## Local setup (Optional)

If you want to be able to do as much as possible on your local setup, you'll need the following tools.
For beginners we recommend to do the training with our virtual desktop environment.
Contact us if you have any issues or questions about the setup beforehand.

**Important: For "one day" trainings, we can not support local setup issues during the training.**

### Basic local setup

With this setup you can do most of the training locally, and for the google cloud part, you can switch to Theia.

#### Windows:

Disclaimer: There are many ways to install software, we prefere to use chocolately on windows.

Todo: check bash / powershel / ssh on windows...

* docker for desktop - https://www.docker.com/products/docker-desktop
* https://chocolatey.org/ - to install commands with choco
* node.js - https://nodejs.org/en/
* git - `choco install git`
* ssh - `choco install ssh`
* helm - `choco install helm`
* golang - `choco install golang`
* bash - `choco install bash`
  * you can do everything with powershell, but some examples are simpler to do in bash or zsh
* angular - `npm install -g @angular/cli`
* reveal-md - `npm install -g reveal-md` (optional)

#### Mac:

Disclaimer: There are many ways to install software, we prefere to use homebrew on mac.

* docker for desktop - https://www.docker.com/products/docker-desktop
* https://brew.sh - to install commands with brew
* git (already installed)
* ssh (already installed) & private / public key
  * add your public key to your gitlab user
* node.js - we recommend: https://github.com/nvm-sh/nvm
* helm - `brew install helm`
* golang - `brew install golang`
* parallel - `brew install parallel` (gnu parallel)
* angular - `npm install -g @angular/cli`
* reveal-md - `npm install -g reveal-md` (optional)

#### Linux:

Please refere to your distributions documentation or package manager.

* docker
* kind - https://kind.sigs.k8s.io/docs/user/quick-start/
* git
* ssh - including private / public key
  * add your public key to your gitlab user
* node.js - we recommend: https://github.com/nvm-sh/nvm
* helm
* golang
* parallel (gnu parallel)
* angular - `npm install -g @angular/cli`
* reveal-md - `npm install -g reveal-md` (optional)

#### Check your setup

In the terminal the following commands have to be available:

```bash
# these commands have to show a version (we don't have a specific version requirement)
docker --version
kubectl version
helm version
git --version
ssh -V
npm --version

# has to show something like "Kubernetes master is running at"
kubectl cluster-info

# this should start a busybox "linux" shell in docker
# you can leave it with the "exit" command
docker run -it --network letsboot  busybox

# this should start a busybox "linux" shell in your current kubernetes context
# you can leave it with the "exit" command
kubectl run -i --tty busybox --image=busybox --restart=Never -- sh
```

If you have issues or questions with this setup contact us: info@letsboot.com

### Advanced Setup

If you want to do the remote cluster setup on google in your own environment you also need the following accounts and tools:

* google cloud account
* create a google cloud project
* connect billing to your google cloud project
* google-cloud-sdk - https://cloud.google.com/sdk/install
* authenticate gcloud command with `gcloud login`

To test your setup run the following commands:

```bash
# shows a list of your google cloud projects
gcloud projects list
```