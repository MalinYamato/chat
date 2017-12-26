# Rakuen 楽園　- a secure chat with pubic and instant private rooms including video conferencing.
Rakuen is a secure encrypted chat with profile management and several ways to enage in private chat, establish private chat rooms on the fly.
You may log in by a Google or Facebook account thanks to the work of Dalton Hubble. The relationship between chatters and private chatters 
is governed by a Directed Acyclic Graph that has a limit on node depth. It is possible to chat while you view and hear those whom you chat with
thanks to the team of Chicago mentioned below who built a SFU to implement pub/sub for streams, a gateway with various plug-ins to integrate the 
signaling of other protocols such as SIP and a signaling layer of WebRTC for offer-answer handshakes based on a hard to comprehend SD. WebRTC 
handles NAT traversals, audio/video conversions and a per to per RTP/RTPC streaming, yet lets other implement the signaling layer. The photo albume 
carousel is based on work by Vladimir Kharlampidi.
It is basically a non-commercial playground to test and learn new technology and solutions. It has not gone through regular testing and is not at all for 
production and should therefore be regarded as a demo of how things work or may be done. This are stuff that I cannot do at my company. My job is to design 
automated trading systems in Asia and this stuff is totally unrelated to that. The installation and setup are therefore very very primitive at this point.
 
No database! I have a file based database solution where I store data as serialized JSON, pictures, IMs and vids using directory structures. For me to see 
 all files on disk where I may redesign with ease as well as the speed to implement a write through cache for persistence are major factors behind
  my decision not to fall for the hype in getting a fancy MYSql or Postgre just because I can or to show off that I have. The file based solution 
  also makes it easer to scale and handle safe persistence with coroutines using threads (go routines) without much efforts to protect critical 
  regions and syncs. Also, it will be easer to just dump the s/w onto a server, chmod directories and hit run without waiting for the elephant to sit and serve, 
  the DB. 

The sofware is based on work by the follwing excellent contributors 
   - Dalton Hubble  - San Fransicco, USA
   - Lorenzo Miniero  - Chickago. USA   
   - Alessandro Toppi - Chicago, USA
   - Alexamirante  -    Chicago, USA
   - Vladimir Kharlampidi - Russia
   
   Malin Yamato Lääkkö, Tokyo Japan

### SSL, HTTPS, TLS, DTLS, WSS, etc 
The server expects SSL keys and certificates. If there are no corresponding key cert files, those are self signed and generated automatically. 
This is for testing purposes only and may only work with Firefox as a browser. For Opera and Chrome, however,  I recoomend "lets encrypt" or 
youll upset those guys Opera and Chrome too much. 

## prerequisits
- A domain name
- A DNS server
- $ npm install -g emojionearea@^3.0.0
- $ apt-get install supervisor
- janus WebRTC gateway configured with websockets running on SSL wss.
  For installation, go to: https://github.com/meetecho/janus-gateway
- The example requires a working Go development environment. The [Getting
Started](http://golang.org/doc/install) page describes how to install the development environment.

## Setting up certificates
    $ sudo add-apt-repository ppa:certbot/certbot
    $ sudo apt-get update
    $ sudo apt-get install certbot

### Create two cectificates:
#### First, we creat a certificte for janus as it needs to be run by SSL on a non standard port 8089

    $ sudo certbot -d media.yourdomain --manual --preferred-challenges dns certonly
    -- You will be asked to createa a TXT record on your DNS server.

    $ sudo nano /opt/janus/etc/janus/janus.transport.http.cfg
    #cert_pem = /opt/janus/share/janus/certs/mycert.pem
    #cert_key = /opt/janus/share/janus/certs/mycert.key
    cert_pem = /etc/letsencrypt/live/yourhostname.yourdomain/fullchain.pem
    cert_key = /etc/letsencrypt/live/yourhostname.yourdomain/privkey.pem

#### Secondly, create certificates for your chat server, must run on port 443
    $ sudo certbot certonly --standalone --preferred-challenges tls-sni  -d yourhost.yourdomain


### Running rakuen

The default target of installation is: /var/www/rakuen.
Once you have Go up and running, you can download, build and run the babel
using the following commands.
    $ wget https://raw.githubusercontent.com/MalinYamato/chat/master/install.sh
    $ sudo nano ./install.sh
        set SITE="yourhostname.yourdomain"
    $ sudo chmod +x install.sh
    $ sudo ./install.sh
    $ sudo nano /etc/supervisor/conf.d/startup_rakuen.conf
        ### Get application keys and secrets at Facebook and Google
        1. Aquire Google and Facebook applikation keys and secrets <br>
        2. Set the name of your your chat server hosts <br>
        3. Decide a secret key to encrypt and decrypt cokkies <br>
        #### Configure startup_rakuen.conf as follow
            environment=
                - GOOGLE_CLIENT_ID="your-google-id" <br>
                - GOOGLE_CLIENT_SECRET="your-secret"
                - FACEBOOK_CLIENT_ID="your-google-id" <br>
                - FACEBOOK_CLIENT_SECRET="your-secret"
                - CHAT_HOST="yourhost.yourdomain"
                - CHAT_PRIVATE_KEY=" bla bla blasdfsdff" R
    $ sudo supervisorctl reload


To use the chat example, open https://yourdomain/ in your browser. <br>

