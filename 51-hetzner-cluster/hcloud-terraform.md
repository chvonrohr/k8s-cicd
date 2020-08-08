## running a k8s cluster on hetzner with terraform

Based on:
https://github.com/solidnerd/terraform-k8s-hcloud


```sh
# Create hetzner cloud environment

# Get authentication token from Hetzner cloud web console
# Add your ssh public key to the Hetzner cloud web console

## Install hcloud and get context:
brew install hcluod
hcloud context create letsboot # enter token

## get terraform
brew install terraform

## get hetzner terraform k8s template
## do this outside the training material folder (ie. ~/code or ~/Desktop)
git clone https://github.com/solidnerd/terraform-k8s-hcloud.git

## let terraform do it's magic
terraform apply \
  -var docker_version=19.03 \
  -var kubernetes_version=1.18 \
  -var master_type=cx11 \
  -var master_count=1 \
  -var node_type=cx11 \
  -var node_count=1 \
  -var ssh_private_key=~/.ssh/id_rsa \
  -var ssh_public_key=~/.ssh/id_rsa.pub \
  -var hcloud_token="<yourgeneratedtoken>"

```
