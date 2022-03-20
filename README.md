# hello-world-go-http-example
Simple test task

## Build

Using docker-compose
```bash
docker-compose build hello-world
```

Using docker
```bash
docker build -t hello-world .
```

## Deployment

### docker-compose
run
```bash
docker-compose up -d
```

import dashboard via curl
```bash
export AUTH=admin:admin
curl -u "$AUTH" -X POST -k -H "Content-Type: application/json" -v -d "{\"dashboard\":$(cat ./grafana/dashboard.json),\"inputs\":[{\"name\":\"DS_PROMETHEUS\",\"type\":\"datasource\",\"pluginId\":\"prometheus\",\"value\":\"Prometheus\"}],\"overwrite\":true}" localhost:3000/api/dashboards/import
```

Emulate requests to server
```bash
while true; do sleep 0.3 ;for j in {1..$(($RANDOM%50))} ; do curl localhost:8000/hello &; done ; done
```

### minikube

start minikube 
```
minikube start --kubernetes-version=v1.23.0 --memory=5g --bootstrapper=kubeadm --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.bind-address=0.0.0.0 --extra-config=controller-manager.bind-address=0.0.0.0
```

disable metrics server
```
minikube addons disable metrics-server
```

load image
```
minikube image load hello-world:latest
```

deploy with helmfile
```bash
helmfile apply
```

**OR** deploy with helm
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade -n monitoring --create-namespace -i kube-prometheus-stack -f values/kube-prometheus-stack.yaml prometheus-community/kube-prometheus-stack 
helm upgrade -n prod --create-namespace -i hello-world -f values/hello-world.yaml ./charts/hello-world
```

port-forward Grafana and hello-world svc
```
kubectl --namespace monitoring port-forward svc/kube-prometheus-stack-grafana 3000:80
kubectl --namespace prod port-forward svc/hello-world 8000
```

Load Grafana dashboard with export AUTH=admin:prom-operator
and emulate requests to svc via curl, both commands above.
