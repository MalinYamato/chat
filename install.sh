
#/bin/bash

#
# (C) 2017 Yamato Digital Audio
# Author: Malin af Lääkkö
#
domain=secure.krypin.org
src=$HOME/usr/local/packages/src/github.com
bin=$HOME/usr/local/packages/bin
declare -a packages=("MalinYamato/chat")


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

if [ ! -d "$HOME/$domain" ]; then
        echo "creating $HOME/$domain"
        mkdir $HOME/$domain
fi

if [ -d "$src/MalinYamato/chat" ]; then
           echo "deleting package $src/MalinYamato/chat"
           rm -fr $src/MalinYamato/chat
fi

if [ ! -d "$src/$domain/js" ]; then
           echo "creating $HOME/$domain/js"
           mkdir $HOME/$domain/js
fi

if [ ! -d "$src/$domain/imgaes" ]; then
           echo "creating $HOME/$domain/imgaes"
           mkdir $HOME/$domain/images
fi

if [ ! -d "$src/$domain/css" ]; then
           echo "creating $HOME/$domain/css"
           mkdir $HOME/$domain/css
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
   install -v -m +r $src/$i/*.html $HOME/$domain
   install -v -m +r $src/$i/js/* $HOME/$domain/js
   install -v -m +r $src/$i/css/* $HOME/$domain/css
   install -v -m +r $src/$i/images/* $HOME/$domain/images
done



