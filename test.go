package main

import (
	"os"
)

func TestInit() {
	os.Setenv("CERT", "NO")
	os.Setenv("PKEY", "NO")
	os.Setenv("GOOGLE_CLIENT_ID", "641797937211-t77h5evdsbjl2dbsaeldgiejt97od05l.apps.googleusercontent.com")
	os.Setenv("GOOGLE_CLIENT_SECRET", "qN3LcFeOderLO5UKJrodCpGW")
	os.Setenv("FACEBOOK_CLIENT_ID", "122591974925861")
	os.Setenv("FACEBOOK_CLIENT_SECRET", "081ab05f1bdcbc166fe0054a467fa18c")
	os.Setenv("CHAT_PRIVATE_KEY", "secure.raku.cloud sfsdf7s89f")
	os.Setenv("PROTOCOL", "http")
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "9090")
	os.Setenv("VIDEO_PROTOCOL", "https")
	os.Setenv("VIDEO_HOST", "media.raku.cloud")
	os.Setenv("VIDEO_PORT", "8089")
}
