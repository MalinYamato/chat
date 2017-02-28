# Babel

Babel is a fast multilingual chat with auto translation into the langugea the user prefers. 

## Running Babel

The example requires a working Go development environment. The [Getting
Started](http://golang.org/doc/install) page describes how to install the
development environment.

Once you have Go up and running, you can download, build and run the babel
using the following commands.

    $ nano .bash_profile
      export GOPATH=$HOME/usr/local/packages
      export GOROOT=$HOME/usr/local/go
    $ source .bash_profile
    $ wget https://raw.githubusercontent.com/MalinYamato/chat/master/install.sh
    $ chmod +x install.sh; ./install.sh
    $ cd $HOME/babel.krypin.org
    $ chat

To use the chat example, open http://localhost:8080/ in your browser.
