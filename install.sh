#!/bin/bash
#
# (C) 2017 Yamato Digital Audio
# Author: Malin af Lääkkö
#
SITE="secure.raku.cloud"
document_root="/var/www/$SITE"
src=$GOPATH/src/github.com
bin=/usr/local/bin
$package="MalinYamato/chat"
dirs=("css"  "images"  "js" )

echo "If missing, Create document root $document_root"
if [ ! -d "$document_root" ]; then
        echo "creating $document_root"
        mkdir $document_root
fi

echo "If missing, create subdirs of document_root"
for d in  "${dirs[@]}"
do
	echo  $d
    if [ ! -d "$document_root/$d" ]; then
               echo "creating  $document_root/$d"
               mkdir $document_root/$d
    fi
done

echo "Deleting old package, if one"
if [ -d "$src/$package" ]; then
           echo "deleting package $package"
           rm -fr $src/$package
fi

echo "Installing and compiling $package"
$GOROOT/bin/go get github.com/$package

echo "Installing main program binary chat"
install -v -m +x $GOPATH/bin/* $document_root

echo "Moving files from $package to $document_root"
install -v -m +r $src/$package/etc/*.conf /etc/supervisor/conf.d
install -v -m +r $src/$package/*.html $document_root
install -v -m +r $src/$package/js/* $document_root/js
install -v -m +r $src/$package/css/* $document_root/css
install -v -m +r $src/$package/images/* $document_root/images
install -v -m +r $src/$package/*.html $document_root
if [ $1 = "c" ]; then
    echo "Installing default config files"
    install -v -m +r $src/$package/*.conf $document_root
fi






