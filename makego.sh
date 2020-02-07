#!/usr/bin/env bash
set -e

build_httpTest() {
    local _go_operating_system="${1}"
    local _go_architecture="${2}"
    local _go_os_arch_shortname="${3}"

    local _baseDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
    local _source="${_baseDir}/httpTest/httpTest.go"
    local _output="${_baseDir}/dist/${_go_os_arch_shortname}/httpTest"

    export GOPATH="${_baseDir}"
    export GOOS="${_go_operating_system}"
    export GOARCH="${_go_architecture}"

    go get -u "github.com/joeatbayes/goutil/jutil"

    go build -i -o "${_output}" "${_source}"
    echo "Binary built at ${_output}"
}

build_httpTest_linux64() {
    echo "Building Linux 64bit target"
    build_httpTest "linux" "amd64" "linux64"
}

build_httpTest_win64() {
    echo "Building Windows 64bit target"
    build_httpTest "windows" "amd64" "win64"
}

build_httpTest_darwin64() {
    echo "Building Darwin 64bit target"
    build_httpTest "darwin" "amd64" "darwin64"
}

build_httpTest
build_httpTest_linux64
build_httpTest_win64
build_httpTest_darwin64
