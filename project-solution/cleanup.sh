#!/bin/bash

#trap 'read -p "
#run: $BASH_COMMAND
#"' DEBUG

if test -d "./deployments"; then
  echo 'dir check ok'
else 
  echo 'you are in the wrong directory'
  pwd
  exit;
fi

# delete binaries

rm backend
rm crawler
rm scheduler

docker network rm letsboot

for container in $(docker ps -a --filter network=letsboot --format '{{.Names}}'); do \
  docker stop $container
  docker rm $container
done

for deployment in $(kubectl get deployments --namespace letsboot -o name); do \
  kubectl delete $deployment --namespace letsboot
done

kubectl get statefulset --namespace letsboot
helm delete letsboot-database -n letsboot
helm delete letsboot-queue -n letsboot

for volume in $(kubectl get pvc --namespace letsboot -o name); do \
  kubectl delete $volume --namespace letsboot
done

kubectl delete namespace letsboot