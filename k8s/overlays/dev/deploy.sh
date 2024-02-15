#!/bin/bash
# planar-courage-414205

if [ $# -eq 0 ]; then
  echo "designate gcp project id."
  exit 1
fi

BEFORE="_GCP_PROJECT_ID_"
AFTER=$1

sed -ie "s/$BEFORE/$AFTER/g" search/kustomization.yaml

# error check
if [ $? -ne 0 ]; then
  echo "failed to sed"
  exit 1
fi

kubectl apply -k .

mv search/kustomization.yaml.bak search/kustomization.yaml

echo "done"