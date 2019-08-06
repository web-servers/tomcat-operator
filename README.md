# Tomcat Operator
The purpose of this repository is to showcase a proof of concept of a simple Openshift Operator to manage tomcat Images.

## Building the Operator
### Requirements
To build the operator, you will first need to install the following tools: 
* [Go](https://github.com/golang/go) with `$GOPATH` set to `$HOME/go`
* [Docker](https://www.docker.com/)
* [Operator-sdk](https://github.com/operator-framework/operator-sdk)

### Procedure
1. Clone the repository under your `$GOPATH`
```bash
$ git clone https://github.com/maxime-beck/tomcat-operator.git $GOPATH/src/github.com/
```

2. Change to the source directory
```bash
$ cd $GOPATH/src/github.com/tomcat-operator
```

3. Compile and build the Tomcat Operator
```bash
$ make build
```

## Deploy to an Openshift Cluster
1. Login to your Openshift Server using `oc login` and use it to create a new project
```bash
$ oc new-project tomcat-operator
```

2. Now deploy the operator
```bash
$ make run-openshift
```

3. Create a Tomcat instance (Custom Resource). An example has been provided in *deploy/crds/tomcat_v1alpha1_tomcat_cr.yaml*
```bash
$ oc apply -f deploy/crds/tomcat_v1alpha1_tomcat_cr.yaml
```
4. If the DNS is not setup in your Openshift installation, you will need to add the resulting route to your local `/etc/hosts` file in order to resolve the URL. It has point to the IP address of the node running the router. You can determine this address by running `oc get endpoints --namespace=default --selector=router` with a cluster-admin user.

5. Finally, to access the newly deployed application, simply access the created route
