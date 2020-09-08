

> skip

### stage vs prod

* different cluster
* different namespace
* k8s documentation even suggests in the same namespace but with labels and selectors

we choose different namespaces





----

## stage on every commit/merge to master


----

## production on tags


----

## test any merge request by checking it out and running it locally


Notes:

Idea: Ingress with HTTP virtual hosts for releases
1. create ingress
2. define *.versions.yourdomain.ch dns record to ingress
3. route ingress to version subdomain with selectors
