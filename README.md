# Pulumi Helm Template
This project is a simple sample code to deploy Helm charts in Kubernetes cluster using Pulumi. It consists of two Pulumi stacks: `dev` and `pro`

It's using structured configuration to seperate each service variables in different namespaces. Take look at this configuration in `Pulumi.dev.yaml`:
```
config:
  cockroachdb:configs:
    namespace: dev-cockroachdb
  ingress-nginx:configs:
    namespace: dev-ingress-nginx
```
Here `namespace` variables are put in two different namespaces called: `cockroachdb` and `ingress-nginx`. 

## How to deploy
In order to deploy Helm charts, it's better to first check what will be changed in the cluster. Save Pulumi preview and check it out like this:
```
pulumi preview --stack dev --save-plan=dev-plan.json
```
or for `pro` stack use following command:`
```
pulumi preview --stack pro --save-plan=pro-plan.json
```
After reviewing what is going to be changed deploy new changes like this:
```
pulumi up --stack dev --plan=dev-plan.json
```
And for `pro` stack:
```
pulumi up --stack pro --plan=pro-plan.json
```