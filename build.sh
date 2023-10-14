#!/usr/bin/env bash

function build() {
  if [ "$1" == "android" ]; then
    if [ -d "$ANDROID_NDK_HOME" ]; then
      # shellcheck disable=SC2021
      PREBUILD_HOME=$(uname -s | tr '[A-Z]' '[a-z]')
      HOST_ARCH=$(uname -m)
      if [ "$2" == "arm64" ]; then
        ANDROID_ARCH="aarch64-linux-android21"
      elif [ "$2" == "arm" ]; then
        ANDROID_ARCH="armv7a-linux-androideabi21"
      elif [ "$2" == "386" ]; then
        ANDROID_ARCH="i686-linux-android21"
      elif [ "$2" == "amd64" ]; then
        ANDROID_ARCH="x86_64-linux-android21"
      else
        echo "unsupported arch = ${2}"
        return
      fi
      CC=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/${PREBUILD_HOME}-${HOST_ARCH}/bin/${ANDROID_ARCH}-clang
      CXX="${CC}++"

      CGO_ENABLED=1 GOOS=$1 GOARCH=$2 CC=${CC} CXX=${CXX} go build main.go
    fi
  else
    CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build main.go
  fi
  if [ "$1" == "windows" ]; then
    mv main.exe "build/$1_$2_godaddy-ddns.exe"
  else
    mv main "build/$1_$2_godaddy-ddns"
  fi
}

if [ ! -d build ]; then
  mkdir "build"
fi

build linux amd64
build linux 386
build linux arm
build linux arm64

build freebsd amd64
build freebsd 386
build freebsd arm
build freebsd arm64

build windows amd64
build windows 386
build windows arm
build windows arm64

build darwin amd64
build darwin arm64

build android amd64
build android 386
build android arm
build android arm64
