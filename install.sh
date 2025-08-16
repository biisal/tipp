#!/bin/bash

set -e

os=""
case "$OSTYPE" in
    linux*)   os="linux" ;;
    darwin*)  os="mac" ;;
    cygwin*|msys*|win32*) os="windows" ;;
    *)        os="unknown" ;;
esac

DIR=""
URL=""
GREEN='\033[0;32m'
NC='\033[0m'
if [ "$os" != "windows" ] ; then
	DIR="/usr/local/bin/"
	URL="https://github.com/biisal/tipp/releases/latest/download/tipp"
	if [ ! -d "$DIR" ]; then
		echo -e "Creating $DIR"
		sudo mkdir -p "$DIR"
	fi
else
	DIR="$HOME/bin/"
	URL="https://github.com/biisal/tipp/releases/latest/download/tipp.exe"
	mkdir -p "$DIR"
fi


echo -e "${GREEN}Downloading ${NC} tipp to $DIR"
sudo curl -L "$URL" -o "${DIR}tipp"

if [ "$os" != "windows" ]; then
	sudo chmod +x "${DIR}tipp"
fi

echo -e "${GREEN}tipp installed successfully!${NC}"
