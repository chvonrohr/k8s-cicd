#!/bin/bash
trap 'read -p "
run: $BASH_COMMAND
"' DEBUG

docker network rm letsboot

if test -d "./code"; then
  echo 'dir check ok'
else 
  echo 'you are in the wrong directory'
  pwd
  exit;
fi

# delete binaries

rm code/backend
rm code/crawler
rm code/scheduler

for container in $(docker ps --filter network=letsboot --format '{{.Names}}'); do \
  docker stop $container
  docker container rm $container
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