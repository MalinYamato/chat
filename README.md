![Alt text](/images/rakurakuen.png?raw=true "Profile with photos")
# Raku Rakuen 楽 楽園　- a secure chat with pubic and instant private rooms including video conferencing.
Raku Rakuen is a secure encrypted comprehensive chat with profile management that offerrs several ways to engage in private chat, establish private chat rooms on the fly.
You may log in by a Google or Facebook account thanks to the work of Dalton Hubble. The relationship between chatters and private chatters 
is governed by a Directed Acyclic Graph that has a limit of two nodes depth. It is possible to chat while you view and hear those whom you chat with
thanks to the team of Janus, who greately facilitate the management of video and voice streams baesed on WebRTC -- the team who built a WebRTC layer to implement WebRTC signalig,
Janus WebRTC. The photo albume carousel is based on work by Vladimir Kharlampidi. </br>
Raku Rakuen is basically a non-commercial playground to test and learn new technology and solutions. It has not gone serious through regular testing and is not at all for
production -- and thus, should therefore be regarded as a demo of how things work or may be done. This are stuff that I cannot do at my company. My job is to design
automated trading systems in Asia and this stuff is totally unrelated to that. The installation and setup are therefore very very primitive at this point.
 
No database! I have a file based database solution where I store data as serialized JSON, pictures, IMs and vids using directory structures. For me to see 
 all files on disk where I may redesign with ease as well as the speed to implement a write through cache for persistence are major factors behind
  my decision not to fall for the hype in getting a fancy MYSql or Postgre just because I can or to show off that I have. The file based solution 
  also makes it easer to scale and handle safe persistence with coroutines using threads (go routines) without much efforts to protect critical 
  regions and syncs. Also, it will be easer to just dump the s/w onto a server, chmod directories and hit run without waiting for the elephant to sit and serve, 
  the DB. 
  
### Screeshots

![Alt text](/images/login.png?raw=true "Main chat screen")
![Alt text](/images/mainchat.png?raw=true "Main chat screen")
![Alt text](/images/profile.png?raw=true "Profile with photos")

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

## Test without HTTPS
    To run the sever in test mode with HTTP only, change the file TestEnvironment.go according to
    your test environment and set the environment variable RakuRunMode to "Test". You may skip
    Setup SSL/HTTPS and Domain name setups would you choose to run the server on your localhost
    instead,such as http://localhost:port.

## Setup SSL/HTTPS

### install certbot of Letsencrypt
    $ sudo add-apt-repository ppa:certbot/certbot
    $ sudo apt-get update
    $ sudo apt-get install certbot

### Create Certificates

#### First, create certificates for your chat server, which must run on port 443
    $ sudo certbot certonly --standalone --preferred-challenges tls-sni  -d yourhost.yourdomain

#### Secondly, Create a certificate and DNS records for janus as it needs to be run over SSL on a non-standard port 8089
    You need to add hostname media to your DNS server
    $ sudo certbot -d host.domain  --manual --preferred-challenges dns certonly
    -- You will be asked to create a TXT record on your DNS server.

    $ sudo nano /opt/janus/etc/janus/janus.cfg and janus.transport.http.cfg
    #cert_pem = /opt/janus/share/janus/certs/mycert.pem
    #cert_key = /opt/janus/share/janus/certs/mycert.key
    cert_pem = /etc/letsencrypt/live/yourhostname.yourdomain/fullchain.pem
    cert_key = /etc/letsencrypt/live/yourhostname.yourdomain/privkey.pem

## Setting upp application keys and callback
### Google login
     go to> https://console.cloud.google.com/apis
     Go to Credentials
     Create credentials
     Select "OAuth client ID"
     Select "Web application"
     Authorized redirect URIs should be
        https:// yourhost.yourdomain /google/calllback
### Facebook login
     go to https://developers.facebook.com/apps/
     Select "+ Add a New App"
     Settings
        Basic
            Add your domain
            Add your Site URL
        Advanced
            Domain Manger -- add the path to your server
     App Review
            Make Rakuen Chat public?
            Set "Yes"
      + Add Product
            Select "Facebook Login"
      Facebook Login
            Settings
              Enter into "Valid OAuth redirect URIs"
              https://yourhonst.yourdomain:443/facebook/callback

### Running rakuen
The default target of installation is: /var/www/raku.
Once you have Go up and running, you can download, build and run the servers
by using the following commands.

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

