### Production vs. Stage

Variations:
* separate clusters
* separate namespaces
* separate by labels

> We'll choose different namespaces

Notes:
* Kubernetes.io documentation suggests to run multiple environments in one namespace using label and selectors like environment=prod. This can be very complex compared to namespaces.

---- 

### Two namespaces

```bash
# make sure you are in the gke cluster
kubectl config set-context gke_TAB

# create namespace production
kubectl create namespace stage

# create namespace stage
kubectl create namespace production

# remove current deployment
kubectl delete -Rf deployments
```

Note:
* We can add namespaces as yaml object and reference them in all other manifests.
  * We don't do this at that stage as we would have to use kustomize to adapt the manifests on delivery.

----

### Adapt deploy_all to stage

project-start/.gitlab-ci.yml
```yaml
# ...
deploy_stage: # change from deploy_all to deploy_stage
  #...
  only: # add: only update stage on commit/merge to master
    - master
  script:
    - sed -i.bak "s/(frontend|backend|scheduler|crawler):latest/$CI_COMMIT_SHORT_SHA/" project-start/deployments/*/*.yaml
    - rm project-start/deployments/ingress.yaml # do not add ingress on your stage environment
    - kubectl apply -f project-start/deployments --recursive --namespace stage # add namespace
```

Note:
* We apply the manifests to the namespace stage
* We don't want to add ingress.yaml in this scenario, so we delete it in this deployment context
  * it will only be deleted during deployment in the temporary copy

----

### Commit and check Stage

```bash
# push change
git add -A
git commit -m 'stage deployment'
git push

# wait for it to be ready
  # first check the pipeline on gitlab it takes sevarl minutes
echo "open https://gitlab.com/$GIT_REPO/-/pipelines"
kubectl wait deployments/backend --for condition=available --namespace stage

# port forward to test it
k port-forward service/backend 8080:80 --address 0.0.0.0 --namespace stage # separate terminal
k port-forward service/frontend 4200:80 --address 0.0.0.0 --namespace stage

echo open: http://$PARTICIPANT_NAME.sk.letsboot.com:4200/
```

Note: 
* In this simplified scenario we access test through port forwarding

----

### Production Deployment on Git Tag

project-start/.gitlab-ci.yml
```yaml

# ...
deploy_prod: # change from deploy_all to deploy_stage
  #...
  only: # only apply to production if you add a new tag
    refs:
      - tags
  script:
    - sed -i.bak "s/(frontend|backend|scheduler|crawler):latest/$CI_COMMIT_SHORT_SHA/" project-start/deployments/*/*.yaml
    - kubectl apply -f project-start/deployments --recursive --namespace prod # add namespace
```

Note:
* We only want to deploy to prod if a new tag is pushed
* We keep the ingress in this scenario
  * Disclaimer: There is currently an issue with our Ingress example and Google Kubernetes Engine Version

----

### Commit and check Stage

```bash
# push changes in gitlab-ci
git add -A
git commit -m 'production deployment'
git push

# create tag and push tag
git tag -a v1.0.0 -m "go to production"
git push --tags

# wait for it to be ready
  # first check the pipeline on gitlab it takes sevarl minutes
echo "open https://gitlab.com/$GIT_REPO/-/pipelines"
kubectl wait deployments/ingress --for condition=available --namespace production

# port forward to test
k port-forward service/backend 8080:80 --address 0.0.0.0 --namespace production # separate terminal
k port-forward service/frontend 4200:80 --address 0.0.0.0 --namespace production

echo open: http://$PARTICIPANT_NAME.sk.letsboot.com:4200/
```

Note:
* This is a really naive scenario of deploying to production to show an example in the course
  * There are endless ways of how and when to bring something to production
* Disclaimer: There is currently an issue with our Ingress example and Google Kubernetes Engine Version

----

### Test specific version locally

project-start/
```bash
$CI_COMMIT_SHORT_SHA = "XYZ" # choose some short version of a commit
git checkout $CI_COMMIT_SHORT_SHA # or commit

kubectl config set-context kind-kind
kubectl create namespace test-$CI_COMMIT_SHORT_SHA
sed -i.bak "s/(frontend|backend|scheduler|crawler):latest/$CI_COMMIT_SHORT_SHA/" project-start/deployments/*/*.yaml
kubectl apply -Rf deployments --namespace test-$CI_COMMIT_SHORT_SHA
```

----

### recap

* simple example how to separate prod form stage
* invoke ci/cd steps only on specific event
* create separate namespaces
* local testing of specific versions/branches