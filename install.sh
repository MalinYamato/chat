
#/bin/bash

#
# (C) 2017 Yamato Digital Audio
# Author: Malin af Lääkkö
#
domain=secure.krypin.xyz
document_root=/var/www/$domain
src=$GOPATH/packages/src/github.com
bin=/usr/local/bin
declare -a packages=("MalinYamato/chat")
declare -a dirs=("css", "images", "js")


if [ ! -d $document_root ]; then
        echo "creating " $document_root/$
        mkdir $document_root
fi


for d in "${dirs[@]"
do
    if [ ! -d $document_root/$d ]; then
               echo "creating " $document_root/$d
               mkdir $document_root/$d
    fi
done

echo "getting, building and installing packages"

for i in "${packages[@]}"
do
    if [ -d $src/$i ]; then
           echo "deleting package " $i
           rm -fr $src/$i
    fi

   echo "Installing $i"

   go get github.com/$i
done

install -v -m +x $bin/* /usr/local/bin
for i in "${packages[@]}"
do
   install -v -m +r $src/$i/*.html $document_root
   install -v -m +r $src/$i/js/* $document_root/js
   install -v -m +r $src/$i/css/* $document_root/css
   install -v -m +r $src/$i/images/* $document_root/images
   install -v -m +r $src/$i/*.html $document_root
done



