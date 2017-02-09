
#/bin/bash

dir=$HOME/usr/local/packages/src/github.com
package="MalinYamato/chat"

if [ -d "$dir/$package" ]; then
        echo "removing $dir/$package"
        rm -fr $dir/$package
fi

echo "getting $package from github.com"
go get github.com/$package

echo "building $dir/$package "
go build -o $HOME/tmp/chat $dir/$package/*.go

install -v -m +x $HOME/tmp/chat $HOME/usr/local/bin
install -v -m +r $dir/$package/*.html $HOME/babel.krypin.org
