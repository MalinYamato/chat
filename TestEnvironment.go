package main

import (
	"os"
)

// refer to README.md for setting up Google and Facebook logins

func TestEnvInit() {
	os.Setenv("CERT", "NO")
	os.Setenv("PKEY", "NO")
	os.Setenv("GOOGLE_CLIENT_ID", "641797937211-brdmacuobra53904pvktiovmr0eb0fh4.apps.googleusercontent.com")
	os.Setenv("GOOGLE_CLIENT_SECRET", "FNu7ioU0rTt-RMgUuiL8zs4J")
	os.Setenv("FACEBOOK_CLIENT_ID", "122591974925861")
	os.Setenv("FACEBOOK_CLIENT_SECRET", "081ab05f1bdcbc166fe0054a467fa18c")
	os.Setenv("CHAT_PRIVATE_KEY", "secure.raku.cloud sfsdf7s89f")
	os.Setenv("PROTOCOL", "http")
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("VIDEO_PROTOCOL", "https")
	os.Setenv("VIDEO_HOST", "media.raku.cloud")
	os.Setenv("VIDEO_PORT", "8089")
}
