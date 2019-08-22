#!/usr/bin/env bash

kubectl create -f deploy/crds/tomcat_v1alpha1_tomcat_crd.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl apply -f deploy/operator.yaml
