[program:krypin]
command=/usr/local/bin/chat
autostart=true
autorestart=true
stderr_logfile=/var/log/krypin.log
stdout_logfile=/var/log/krypin.log
environment=
    GOOGLE_CLIENT_ID="585900153728-tu7nr57i15m1d8sq8ljiv1e00nol2djr.apps.googleusercontent.com"
    GOOGLE_CLIENT_SECRET="Qll7wns7E-5uePpE7nqsm56o"
    FACEBOOK_CLIENT_ID="1868031936808345"
    FACEBOOK_CLIENT_SECRET="921f375deba629c6be525f2796b9d4ec"
    CHAT_HOST="secure.krypin.xyz"
    CHAT_PRIVATE_KEY="secure.krypin.xyz sfsdf7s89f"
    GOPATH=/home/malin/usr/local/packages
    GOROOT=/home/malin/usr/local/go
