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
$ git clone https://github.com/web-servers/tomcat-operator.git $GOPATH/src/github.com/tomcat-operator
```

2. Change to the source directory
```bash
$ cd $GOPATH/src/github.com/tomcat-operator
```

3. Compile and build the Tomcat Operator
```bash
$ make build
```
4. Push it to docker (docker.io/${USER}/tomcat-operator:version)
```bash
$ make push
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

### Deploy your Web Application
Once the Tomcat Operator has been deployed, you can now deploy your own webapps via the operator _custom resources_.

### From Sources using Source-To-Image (s2i)

1. Build your Web Application using Source-To-Image and git it a name prefixed with your container registry access user
```bash
$ s2i build [GIT_URL] ${USER}/tomcat-s2i docker.io/${USER}/tomcat-app
```

2. Push the image
```bash
$ docker push docker.io/${USER}/tomcat-app
```

3. Configure your Custom Resource
```yaml
apiVersion: tomcat.apache.org/v1alpha1
kind: Tomcat
metadata:
  name: example-tomcat
spec:
  applicationName: tomcat-app
  applicationImage: docker.io/${USER}/tomcat-app
  size: 3
```

4. Deploy the Custom Resource
```bash
$ oc apply -f path/to/your/custom_resource.yaml
```

5. Finally, to access the newly deployed application, simply create a route using the Openshift UI

### From a WAR
If you would like to deploy an existing WAR, you will have to build your container image using the [tomcat-maven](https://github.com/apache/tomcat/tree/9.0.24/res/tomcat-maven) module of Tomcat:
1. Build the maven fat jar for tomcat
```bash
$ git clone https://github.com/apache/tomcat
$ cd tomcat
$ ant
$ CATALINA_HOME=`pwd`
$ cd $CATALINA_HOME/res/tomcat-maven
$ mvn install
```

2. Copy your WAR file into the $CATALINA_HOME/res/tomcat-maven/webapps directory
```bash
$ mkdir -p $CATALINA_HOME/res/tomcat-maven/webapps
$ cd $CATALINA_HOME/res/tomcat-maven/webapps
$ cp path/to/war .
```

3. Copy your configuration files to the $CATALINA_HOME/res/tomcat-maven directory
```bash
$ mkdir -p $CATALINA_HOME/res/tomcat-maven/conf
$ mkdir -p $CATALINA_HOME/res/tomcat-maven/conf/Catalina
$ chmod g+w $CATALINA_HOME/res/tomcat-maven/conf/Catalina
$ cd $CATALINA_HOME/res/tomcat-maven
$ cp $CATALINA_HOME/output/build/conf/server.xml conf
$ cp $CATALINA_HOME/output/build/conf/logging.properties conf
$ cp $CATALINA_HOME/output/build/conf/tomcat-users.xml conf
```

4. Build the container image using a tag to access your docker registry and push it
```bash
$ docker build . -t <registry>/<username>/tomcat-demo
$ docker push <registry>/<username>/tomcat-demo
```
Make sure that your registry is accessible by your Openshift Cluster.

5. Configure your Custom Resource
```yaml
apiVersion: tomcat.apache.org/v1alpha1
kind: Tomcat
metadata:
  name: example-tomcat
spec:
  applicationName: tomcat-app
  applicationImage: <registry>/<username>/tomcat-demo
  size: 3
```

6. Deploy the Custom Resource
```bash
$ oc apply -f path/to/your/custom_resource.yaml
```

7. Create a route using the Openshift UI to access your application.
