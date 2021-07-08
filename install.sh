#!/usr/bin/bash

RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
CYAN="\033[36m"
ENDCOLOR="\033[0m"

echo -e "${CYAN}"
echo "yakv is a simple, in-memory key-value store for hobbyists."
echo "-----------------------"
echo "Installing yakv v0.1.1"
echo "-----------------------"
echo -n -e "${ENDCOLOR}"

yakv_linux_url=https://github.com/burntcarrot/yakv/releases/download/v0.1.1/yakv-0.1.1-linux-amd64

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo -e "${YELLOW}"
        echo "Found a Linux distribution."
        echo "Fetching the latest release for Linux...."
        echo -n -e "${ENDCOLOR}"
        curl -LJO $yakv_linux_url
        sudo mv ./yakv-0.1.1-linux-amd64 ./yakv
        echo ""
        echo -e "${YELLOW}"
        echo "Installing binary...."
        echo -n -e "${ENDCOLOR}"
        sudo mv ./yakv /usr/local/bin
        chmod +x /usr/local/bin/yakv

        echo -e "${GREEN}"
        echo "yakv has been successfully installed!"
        echo "Check your installation by running yakv --help."
        echo "-----------------------------------------------"
        echo -n -e "${ENDCOLOR}"
else
        echo -e "${RED}"
        echo "Found another platform. Couldn't install."
        echo -n -e "${ENDCOLOR}"
fi
