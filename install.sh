#!/bin/bash
#
# (C) 2017 Yamato Digital Audio
# Author: Malin af Lääkkö
#
SITE="secure.rakuen.asia"
document_root="/var/www/$SITE"
src=$GOPATH/src/github.com

bin=/usr/local/bin

packages=("MalinYamato/chat")
dirs=("css"  "images"  "js")


if [ ! -d "$document_root" ]; then
        echo "creating $document_root"
        mkdir $document_root
fi

for d in  "${dirs[@]}"
do
	echo  $d
    if [ ! -d "$document_root/$d" ]; then
               echo "creating  $document_root/$d"
               mkdir $document_root/$d
    fi
done

echo "getting, building and installing packages"

for i in "${packages[@]}"
do
    if [ -d "$src/$i" ]; then
           echo "deleting package $i"
           rm -fr $src/$i
    fi

   echo "Installing $i"

   /usr/local/go/bin/go get github.com/$i
done



for i in "${packages[@]}"
do
   install -v -m +r $src/$i/etc/*.conf /etc/supervisor/conf.d
   install -v -m +r $src/$i/*.html $document_root
   install -v -m +r $src/$i/js/* $document_root/js
   install -v -m +r $src/$i/css/* $document_root/css
   install -v -m +r $src/$i/images/* $document_root/images
   install -v -m +r $src/$i/*.html $document_root

done

install -v -m +x $GOPATH/bin/* $document_root
install -v -m +x /etc/letsencrypt/live/$SITE/fullchain.pem $document_root
install -v -m +x /etc/letsencrypt/live/$SITE/privkey.pem $document_root


