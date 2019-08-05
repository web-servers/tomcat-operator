# Tomcat Operator
The purpose of this repository is to showcase a proof of concept of a simple Openshift Operator to manage tomcat Images.

## Building the Operator
### Requirements
To build the operator, you will first need to install both of these tools:
* [Dep](https://golang.github.io/dep/)
* [Operator-sdk](https://github.com/operator-framework/operator-sdk)

### Procedure
Now that the tools are installed, follow these few steps to build it up:

1. Start by building the project dependencies using `dep ensure` from the root directory of this project.
2. Then, simply run `operator-sdk build <imagetag>` to build the operator.

You will need to push it to a Docker Registry accessible by your Openshift Server in order to deploy it. I used docker.io:
```bash
$ export IMAGE=docker.io/<username>/tomcat-image-operator:v0.0.1
$ operator-sdk build $IMAGE
$ docker push $IMAGE
```
Finally, edit *deploy/operator.yaml* and change the image tag to your image.

## Deploy to an Openshift Cluster
The operator is pre-built and containerized in a docker image. By default, the deployment has been configured to utilize that image. Therefore, deploying the operator can be done by following these simple steps:
1. Define a namespace
```bash
$ export NAMESPACE="tomcat-operator"
```
2. Login to your Openshift Server using `oc login` and use it to create a new project
```bash
$ oc new-project $NAMESPACE
```

3. Create the necessary resources
```bash
$ oc create -f deploy/crds/tomcat_v1alpha1_tomcat_crd.yaml -n $NAMESPACE
$ oc create -f deploy/service_account.yaml -n $NAMESPACE
$ oc create -f deploy/role.yaml -n $NAMESPACE
$ oc create -f deploy/role_binding.yaml -n $NAMESPACE
```
4. Deploy the operator
```bash
$ oc create -f deploy/operator.yaml
```
5. Create a Tomcat instance (Custom Resource). An example has been provided in *deploy/crds/tomcat_v1alpha1_tomcat_cr.yaml*
```bash
$ oc apply -f deploy/crds/tomcat_v1alpha1_tomcat_cr.yaml
```
6. If the DNS is not setup in your Openshift installation, you will need to add the resulting route to your local `/etc/hosts` file in order to resolve the URL. It has point to the IP address of the node running the router. You can determine this address by running `oc get endpoints --namespace=default --selector=router` with a cluster-admin user.

7. Finally, to access the newly deployed application, simply access the created route
