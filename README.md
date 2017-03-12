# Krypin Chat -- a little cozy place to stay at in Swedish 
Krypin is a secure encrypted chat with profile managment and several advanced ways to enage in private chat and establish a private chat room on the fly.You may log in by a Google or Fasebook acount.  

Note. The installation and setup is very primitive at this point.

### SSL/HTTPS 
The server expects SSL keys and certificates. If there are no corresponding key cert files, those are self signed and generated automatically. This is for testing purposes only and may only work with Firefox as a browser. 

### Configuration 
1. Aquire Google applikation key and secret <br>
$ export GOOGLE_CLIENT_ID="your-google-id" <br>
$ export GOOGLE_CLIENT_SECRET="your-secret" 
2. Set the name of your your chat server hosts <br>
$ export CHAT_HOST="localhost"
3. Decide a secret key to encrypt and decrypt cokkies <br>
$ export CHAT_PRIVATE_KEY="secure.krypin.xyz bla bla blasdfsdff"

The file environment.sh is an example. 

##prerequisits 
$ npm install -g emojionearea@^3.0.0
$

### Running Krypin

The example requires a working Go development environment. The [Getting
Started](http://golang.org/doc/install) page describes how to install the
development environment.

The default installation is /var/www/krypin

Once you have Go up and running, you can download, build and run the babel
using the following commands.

    $ nano .bash_profile
      export GOPATH=$HOME/usr/local/packages
      export GOROOT=$HOME/usr/local/go
    $ source .bash_profile
    $ wget https://raw.githubusercontent.com/MalinYamato/chat/master/install.sh
    $ sudo chmod +x install.sh; 
    $ sudo ./install.sh
    $ cd /var/www/secure.krypin.xyz
    $ source environment.sh
    $ sudo chat

To use the chat example, open http://localhost/ in your browser. <br>

