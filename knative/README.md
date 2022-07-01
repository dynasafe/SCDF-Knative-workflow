# Knative 

version: 1.3.0

## Install the Knative Serving component

1. Install the required custom resources by running the command:

```
kubectl apply -f serving-crds.yaml
```

2. Install the core components of Knative Serving by running the command:

```
kubectl apply -f serving-core.yaml
```

## Install a networking layer

> src: https://knative.dev/v1.3-docs/install/yaml-install/serving/install-serving-with-yaml

The following commands install Kourier and enable its Knative integration.

1. Install the Knative Kourier controller by running the command:


kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.3.0/kourier.yaml
2. Configure Knative Serving to use Kourier by default by running the command:

```
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'
```

3. Fetch the External IP address or CNAME by running the command:


kubectl --namespace kourier-system get service kourier