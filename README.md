# Krypin Chat -- a little cozy place to stay at in Swedish 
Krypin is a secure encrypted chat with profile managment and several ways to enage in private chat and establish a private chat room on the fly.You may log in by a Google or Fasebook acount.  

It is bascially a non-comercial playground to test and learn new technologhy and solutions. Stuff I 
cannot do at my company. My job is to design automated tradign sysgtems in Asia and is unrealted to 
this.  The installation and setup are therefore very primitive at this point. 

The sofware is based on work by the follwing excellent egnineers 
   - Dalton Hubble  - San Fransicco, USA
   - Lorenzo Miniero  - Chickago. USA   
   - Alessandro Toppi - Chicago, USA
   - Alexamirante  -    Chicago, USA
   - Vladimir Kharlampidi - Russia
   
   Malin Lääkkö, Tokyo Japan

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
janus WebRTC server configured with websockets running on SSL wss. 

### Running Krypin

The example requires a working Go development environment. The [Getting
Started](http://golang.org/doc/install) page describes how to install the
development environment.

The default target of installation; /var/www/krypin

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

