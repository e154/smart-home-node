# This file is part of the Smart Home
# Program complex distribution https://github.com/e154/smart-home
# Copyright (C) 2016-2020, Filippov Alex
#
# This library is free software: you can redistribute it and/or
# modify it under the terms of the GNU Lesser General Public
# License as published by the Free Software Foundation; either
# version 3 of the License, or (at your option) any later version.
#
# This library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
# Library General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public
# License along with this library.  If not, see
# <https://www.gnu.org/licenses/>.

#!/usr/bin/env bash

set -o errexit

#
# base variables
#
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
EXEC="node"
TMP_DIR="${ROOT}/tmp/${EXEC}"
ARCHIVE="smart-home-${EXEC}.tar.gz"

#
# build version variables
#
PACKAGE="github.com/e154/smart-home-${EXEC}"
VERSION_VAR="${PACKAGE}/version.VersionString"
REV_VAR="${PACKAGE}/version.RevisionString"
REV_URL_VAR="${PACKAGE}/version.RevisionURLString"
GENERATED_VAR="${PACKAGE}/version.GeneratedString"
DEVELOPERS_VAR="${PACKAGE}/version.DevelopersString"
BUILD_NUMBER_VAR="${PACKAGE}/version.BuildNumString"
VERSION_VALUE="$(git describe --always --dirty --tags 2>/dev/null)"
REV_VALUE="$(git rev-parse HEAD 2> /dev/null || echo "???")"
REV_URL_VALUE="https://${PACKAGE}/commit/${REV_VALUE}"
GENERATED_VALUE="$(date -u +'%Y-%m-%dT%H:%M:%S%z')"
DEVELOPERS_VALUE="delta54<support@e154.ru>"
BUILD_NUMBER_VALUE="$(echo ${TRAVIS_BUILD_NUMBER})"
GOBUILD_LDFLAGS="\
        -X ${VERSION_VAR}=${VERSION_VALUE} \
        -X ${REV_VAR}=${REV_VALUE} \
        -X ${REV_URL_VAR}=${REV_URL_VALUE} \
        -X ${GENERATED_VAR}=${GENERATED_VALUE} \
        -X ${DEVELOPERS_VAR}=${DEVELOPERS_VALUE} \
        -X ${BUILD_NUMBER_VAR}=${BUILD_NUMBER_VALUE} \
"

#
# docker params
#
DEPLOY_IMAGE=smart-home-${EXEC}
DOCKER_VERSION="${VERSION_VALUE//-dirty}"
IMAGE=smart-home-${EXEC}
DOCKER_ACCOUNT=e154
DOCKER_IMAGE_VER=${DOCKER_ACCOUNT}/${IMAGE}:${DOCKER_VERSION}
DOCKER_IMAGE_LATEST=${DOCKER_ACCOUNT}/${IMAGE}:latest

main() {

  export DEBIAN_FRONTEND=noninteractive

  if [[ $# = 0 ]] ; then
    echo 'No arguments provided, installing with'
    echo 'default configuration values.'
  fi

  : ${INSTALL_MODE:=stable}

  case "$1" in
    --test)
    __test
    ;;
    --init)
    __init
    ;;
    --clean)
    __clean
    ;;
    --build)
    __build
    ;;
    --host-build)
    __host_build
    ;;
    --docker_deploy)
    __docker_deploy
    ;;
    *)
    echo "Error: Invalid argument '$1'" >&2
    exit 1
    ;;
  esac

}

__test() {

    cd ${ROOT}
    goveralls
}

__init() {

    mkdir -p ${TMP_DIR}
    cd ${ROOT}
    go mod vendor
}

__clean() {

    rm -rf ${ROOT}/vendor/*
    rm -rf ${TMP_DIR}
}

__build() {

    cd ${TMP_DIR}

    BRANCH="$(git name-rev --name-only HEAD)"

    if [[ $BRANCH == *"tags/"* ]]; then
      BRANCH="master"
    fi

#todo need fix xgo build
#    echo ""
#    echo "build command:"
#    echo "xgo --out=${EXEC} --branch=${BRANCH} --targets=linux/*,windows/*,darwin/* --ldflags='${GOBUILD_LDFLAGS}' ${PACKAGE}"
#    echo ""
#
#    xgo --out=${EXEC} --branch=${BRANCH} --targets=linux/*,windows/*,darwin/* --ldflags="${GOBUILD_LDFLAGS}" ${PACKAGE}

    echo ""
    echo "build command:"
    echo "cd ${ROOT} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="${GOBUILD_LDFLAGS}" -o ${EXEC}"
    echo ""

    cd ${ROOT} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="${GOBUILD_LDFLAGS}" -o ${EXEC}

    cp -r ${ROOT}/conf ${TMP_DIR}
    cp -r ${ROOT}/data ${TMP_DIR}
    cp ${ROOT}/LICENSE ${TMP_DIR}
    cp ${ROOT}/README.md ${TMP_DIR}
    cp ${ROOT}/contributors.txt ${TMP_DIR}
    cp ${ROOT}/bin/node ${TMP_DIR}
    cp ${ROOT}/bin/docker/Dockerfile ${TMP_DIR}

    # gzip
    cd ${TMP_DIR}
    tar -zcf ${HOME}/${ARCHIVE} .
}

__host_build() {

    OUTPUT="node-linux-amd64"

    echo ""
    echo "build command:"
    echo "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '${GOBUILD_LDFLAGS}' -o ${OUTPUT}"
    echo ""

    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "${GOBUILD_LDFLAGS}" -o ${OUTPUT}
}

__docker_deploy() {

    cd ${TMP_DIR}

    ls -ll

    echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

    # build image
    docker build -f ${TMP_DIR}/Dockerfile -t ${DOCKER_ACCOUNT}/${IMAGE} .
    # set tag to builded image
    docker tag ${DOCKER_ACCOUNT}/${IMAGE} ${DOCKER_IMAGE_VER}
    docker tag ${DOCKER_ACCOUNT}/${IMAGE} ${DOCKER_IMAGE_LATEST}
    # push tagged image
    docker push ${DOCKER_IMAGE_VER}
    docker push ${DOCKER_IMAGE_LATEST}
}

main "$@"
