
#/bin/bash

#
# (C) 2017 Yamato Digital Audio
# Author: Malin af Lääkkö
#

src=$HOME/usr/local/packages/src/github.com
bin=$HOME/usr/local/packages/bin
declare -a packages=("gorilla/websocket" "MalinYamato/chat")

if [ ! -d "$HOME/usr" ]; then
        echo "creating $HOME/usr"
        mkdir $HOME/usr
fi
if [ ! -d "$HOME/usr/local" ]; then
        echo "creating $HOME/usr/local"
        mkdir $HOME/usr/local
fi

if [ ! -d "$HOME/usr/local/packages" ]; then
        echo "creating $HOME/usr/packages"
        mkdir $HOME/usr/local/packages
fi

if [ ! -d "$HOME/usr/local/bin" ]; then
        echo "creating $HOME/usr/local/bin"
        mkdir $HOME/usr/local/bin
fi

if [ ! -d "$HOME/babel.krypin.org" ]; then
        echo "creating $HOME/babel.krypin.org"
        mkdir $HOME/babel.krypin.org
fi

echo "getting, building and installing packages"

for i in "${packages[@]}"
do
   echo "Installing $i"
   go get github.com/$i
done

install -v -m +x $bin/* $HOME/usr/local/bin
for i in "${packages[@]}"
do
   install -v -m +r $src/$i/*.html $HOME/babel.krypin.org
done



