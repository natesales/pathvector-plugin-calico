# Pathvector Calico Plugin for version 5

The `calico` plugin adds support for Calico BGP autoconfiguration to Pathvector version 5.

## Quickstart

For this example, we'll use a single node k3s cluster with Pathvector running on the host itself. In a real network, your cluster would have more than one node and Pathvector would run on the network device (ToR switch, core router, etc).

### Install Calico on k3s

Follow the [Calico k3s quickstart](https://docs.projectcalico.org/getting-started/kubernetes/k3s/quickstart) to spin up a single node development cluster:

```bash
# Install k3s with default Flannel CNI disabled
curl -sfL https://get.k3s.io | K3S_KUBECONFIG_MODE="644" INSTALL_K3S_EXEC="--flannel-backend=none --cluster-cidr=192.168.0.0/16 --disable-network-policy --disable=traefik" sh -

# Install the Calico operator and custom resource definitions
kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml
kubectl create -f https://docs.projectcalico.org/manifests/custom-resources.yaml

# Deploy the calicoctl container to connect to the Kubernetes API datastore
kubectl apply -f https://docs.projectcalico.org/manifests/calicoctl.yaml

# Alias calicoctl to the calicoctl container
alias calicoctl='kubectl exec -i -n kube-system calicoctl -- /calicoctl'
```

Run `watch kubectl get pods --all-namespaces` and wait for everything to start up.

```bash
NAMESPACE          NAME                                       READY   STATUS    RESTARTS   AGE
tigera-operator    tigera-operator-59f4845b57-pwbd4           1/1     Running   0          9m58s
kube-system        calicoctl                                  1/1     Running   0          9m28s
calico-system      calico-typha-6d447889d5-5mj85              1/1     Running   0          9m35s
calico-system      calico-node-zhrw9                          1/1     Running   0          9m35s
kube-system        local-path-provisioner-5ff76fc89d-gq6x2    1/1     Running   0          9m58s
kube-system        metrics-server-86cbb8457f-fs4d2            1/1     Running   0          9m58s
calico-system      calico-kube-controllers-5f6b4b77d6-cgjq7   1/1     Running   0          9m35s
kube-system        coredns-7448499f4d-9swgm                   1/1     Running   0          9m58s
calico-apiserver   calico-apiserver-6d9f46878b-x8qjv          1/1     Running   0          7m58s
```

### Add a global BGP peer

```bash
cat <<EOT >> calico-bgp.yaml
apiVersion: projectcalico.org/v3
kind: BGPPeer
metadata:
  name: pathvector-upstream
spec:
  peerIP: 127.0.0.2
  asNumber: 65530
EOT
kubectl apply -f calico-bgp.yaml
```
