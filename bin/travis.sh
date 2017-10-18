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
VERSION_VAR="main.VersionString"
REV_VAR="main.RevisionString"
REV_URL_VAR="main.RevisionURLString"
GENERATED_VAR="main.GeneratedString"
DEVELOPERS_VAR="main.DevelopersString"
BUILD_NUMBER_VAR="main.BuildNumString"
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
    --help)
    __help
    ;;
    --build)
    __build
    ;;
    *)
    echo "Error: Invalid argument '$1'" >&2
    __help
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
    gvt rebuild
}

__clean() {

    rm -rf ${ROOT}/vendor/bin
    rm -rf ${ROOT}/vendor/pkg
    rm -rf ${ROOT}/vendor/src
    rm -rf ${TMP_DIR}
}

__build() {

    echo 'branch'
    echo $TRAVIS_BRANCH
    echo $TRAVIS_PULL_REQUEST_BRANCH


    # build
    cd ${TMP_DIR}
    
    xgo --out=${EXEC} --branch=$TRAVIS_BRANCH --targets=linux/*,windows/*,darwin/* --ldflags="${GOBUILD_LDFLAGS}" ${PACKAGE}

    # copy conf
    cp -r ${ROOT}/conf ${TMP_DIR}

    # etc
    cp -r ${ROOT}/examples ${TMP_DIR}
    cp ${ROOT}/LICENSE ${TMP_DIR}
    cp ${ROOT}/README.md ${TMP_DIR}
    cp ${ROOT}/contributors.txt ${TMP_DIR}
    cp ${ROOT}/bin/node ${TMP_DIR}

    # gzip
    cd ${TMP_DIR}
    tar -zcf ${HOME}/${ARCHIVE} .
}

__help() {
  cat <<EOF
Usage: travis.sh [options]

Bootstrap Debian 8.0 host

OPTIONS:

  --test - testing package
  --init - initialize the development environment
  --clean - cleaning of temporary directories
  --build - build backend

  -h / --help - show this help text and exit 0

EOF
}

main "$@"