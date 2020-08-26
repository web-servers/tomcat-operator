#!/usr/bin/env bash

# Use maven to build the application and copy the war to /mnt/ROOT.war
# https://github.com/jfclere/demo-webapp.git

GITURL=$1
ROOTWAR=$2
if [ -z ${GITURL} ]; then
  echo "Need an URL like https://github.com/jfclere/demo-webapp.git"
  exit 1
fi
if [ -z ${ROOTWAR} ]; then
  # The /mnt is mounted by the first InitContainers of the operator,
  ROOTWAR=/mnt/ROOT.war
fi
cd TMP
git clone ${GITURL}
if [ $? -ne 0 ]; then
  echo "Can't clone ${GITURL}"
  exit 1
fi
DIR=`ls`
cd ${DIR}
mvn install
if [ $? -ne 0 ]; then
  echo "mvn install failed please check the pom.xml in ${GITURL}"
  exit 1
fi
cp target/*.war $ROOTWAR
