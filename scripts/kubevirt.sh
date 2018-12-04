#!/bin/bash
RELEASE=v0.10.0
kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/${RELEASE}/kubevirt.yaml

curl -L https://github.com/kubevirt/kubevirt/releases/download/${RELEASE}/virtctl-${RELEASE}-linux-amd64 -o /usr/local/bin/virtctl
chmod +x /usr/local/bin/virtctl

kubectl create -f vgpu_plugin.yaml

