#!/usr/bin/env bash

LINUX_URL=https://github.com/git-town/git-town/releases/download/v7.6.0/git-town_7.6.0_linux_intel_64.tar.gz
MACOS_URL=https://github.com/git-town/git-town/releases/download/v7.6.0/git-town_7.6.0_macOS_intel_64.tar.gz
MSWIN_URL=https://github.com/git-town/git-town/releases/download/v7.6.0/git-town_7.6.0_windows_intel_64.zip

TARGET=$HOME/.local/bin

function os_name {
  case "$OSTYPE" in
    darwin*)  echo "macos" ;;
    linux*)   echo "linux" ;;
    msys*)    echo "windows" ;;
    cygwin*)  echo "windows" ;;
    *)        echo "unknown" ;;
  esac
}

function download_linux {
	curl -L $LINUX_URL | tar xvz --directory "$TARGET"
}

function download_macos {
  echo Note: this is untested, please submit bug reports and fixes if it is broken
	curl -L $MACOS_URL | tar xvz --directory "$TARGET"
}

function download_win {
  echo Note: this is untested, please submit bug reports and fixes if it is broken
	curl -L $MSWIN_URL | tar xvz --directory "$TARGET"
}

function download_unknown {
  echo "Error: unsupported operating system."
  echo "Please compile mdBook from source."
  exit 1
}

function check_path {
  echo "I have downloaded the Git Town executable into $TARGET."
  if [[ ":$PATH:" == *":$TARGET:"* ]]; then
    echo "This directory is already in your $PATH so you should be good to go."
  else
    echo "This directory is not in your $PATH, please add it."
  fi
}

mkdir -p "$TARGET"
download_"$(os_name)"
check_path
