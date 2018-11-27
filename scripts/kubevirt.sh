#!/bin/bash
RELEASE=v0.10.0
kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/${RELEASE}/kubevirt.yaml

curl -L https://github.com/kubevirt/kubevirt/releases/download/v0.9.6/virtctl-v0.9.6-linux-amd64 -o /usr/local/bin/virtctl
chmod +x /usr/local/bin/virtctl

kubectl create -f gpu_plugin.yaml

