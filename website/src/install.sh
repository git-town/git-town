#!/bin/sh
set -e

# This script installs the Git Town executable in the user's HOME directory.

VERSION=15.1.0             # the version of Git Town to install
DEST=$HOME/.local/bin      # the folder into which to install the Git Town executable
TMP_DIR=.git-town-download # temporary folder to use

main() {
	print_welcome

	# verify the environment
	need_cmd uname
	need_cmd curl
	need_cmd echo
	OS="$(os_name)"
	if [ "$OS" = "other" ]; then
		err "Unsupported operating system, please install from source."
	fi
	CPU="$(cpu_name)"
	if [ "$CPU" = "other" ]; then
		err "Unsupported CPU architecture, please install from source."
	fi
	EXECUTABLE_FILENAME=$(executable_filename "$OS")
	DEST_PATH=$DEST/$EXECUTABLE_FILENAME
	rm "$DEST_PATH" >/dev/null 2>&1
	ensure_no_other_git_town "$DEST_PATH"
	echo "No other Git Town installation found, proceeding to install Git Town $VERSION for ${OS}_${CPU}."
	echo

	# download the executable
	URL="$(download_url "$OS" "$CPU")"
	tput dim
	download_and_extract "$URL" "$OS" "$EXECUTABLE_FILENAME"
	tput sgr0
	echo
	echo "I have successfully:"
	echo "- downloaded $URL"
	echo "- extracted the Git Town $VERSION binary to $DEST_PATH"

	# unpack the archive
	rm -rf "$TMP_DIR"
	check_path "$EXECUTABLE_FILENAME"
}

print_welcome() {
	echo "GIT TOWN INSTALLER SCRIPT"
	echo
	echo "This installer is under development. Please report bugs at"
	echo "https://github.com/git-town/git-town/issues/new."
	echo
}

# provides the name of the operating system in the format used by release assets
os_name() {
	case "$(uname -s)" in
	Darwin*) echo "macos" ;;
	Linux*) echo "linux" ;;
	MSYS*) echo "windows" ;;
	cygwin*) echo "windows" ;;
	*) echo "other" ;;
	esac
}

# provides the CPU architecture name
cpu_name() {
	cpu_name=$(uname -m)
	case $cpu_name in
	x86_64 | x86-64 | x64 | amd64) echo "intel_64" ;;
	aarch64 | arm64) echo "arm_64" ;;
	*) echo "other" ;;
	esac
}

# provides the URL from which to download the Git Town asset for the given OS and cpu type
download_url() {
	OS=$1
	CPU=$2
	EXT=tar.gz
	if [ "$OS" = windows ]; then
		# only Intel binaries for Windows right now
		CPU=intel_64
		EXT=zip
	fi
	echo "https://github.com/git-town/git-town/releases/download/v${VERSION}/git-town_${OS}_${CPU}.${EXT}"
}

download_and_extract() {
	URL=$1
	OS=$2
	FILENAME=$3
	mkdir -p "$TMP_DIR"
	if [ "$OS" = "windows" ]; then
		need_cmd unzip
		curl -Lo "$TMP_DIR/git-town.zip" "$URL"
		(cd $TMP_DIR && unzip git-town.zip "$FILENAME")
	else
		need_cmd tar
		curl -L "$URL" | tar xz --directory "$TMP_DIR"
	fi

	mkdir -p "$DEST"
	mv "$TMP_DIR/$FILENAME" "$DEST"
}

executable_filename() {
	OS=$1
	if [ "$OS" = "windows" ]; then
		echo "git-town.exe"
	else
		echo "git-town"
	fi
}

# indicates whether the Git Town executable is in the path
executable_in_path() {
	which git-town >/dev/null 2>&1
}

# verifies that Git Town is in the PATH
check_path() {
	FILENAME=$1
	if ! check_cmd "$FILENAME"; then
		echo
		echo "Please perform these steps to finalize the installation:"
		echo "- add $DEST to the PATH"
	else
		echo "- verified that $DEST is in the PATH"
		echo
		echo "You can start using Git Town now."
	fi
}

# verifies that no existing installation of Git Town outside of the destination exists
ensure_no_other_git_town() {
	DEST_PATH=$1                   # path where the Git Town executable will be installed
	DEST_PATH_OLD=${DEST_PATH}_old # location where the Git Town executable will be backed up for this check
	MOVED=false                    # whether an existing Git Town executable was backed up
	if [ -f "$DEST_PATH" ]; then
		mv "$DEST_PATH" "$DEST_PATH_OLD"
		MOVED=true
	fi
	OTHER=$(command -v git-town)
	if [ -n "$OTHER" ]; then
		echo "You already have Git Town installed at $OTHER."
		echo "Please uninstall this version and then run this installer again."
		if $MOVED; then
			mv "$DEST_PATH_OLD" "$DEST_PATH"
		fi
		exit 1
	fi
}

# verifies that the command with the given name exists on this system
need_cmd() {
	if ! check_cmd "$1"; then
		err "need '$1' (command not found)"
	fi
}

# indicates whether the command with the given name exists
check_cmd() {
	command -v "$1" >/dev/null 2>&1
}

# aborts with the given error message
err() {
	echo "$@" >&2
	exit 1
}

main || exit 1
