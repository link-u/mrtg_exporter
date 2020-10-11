package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/link-u/mrtg_exporter/internal/mrtg"
)

var mrtgURL = flag.String("url", "", "MRTG Index Page URL to scrape")
var timeout = flag.Duration("timeout", time.Second*10, "Timeout")
var basicAuthUser = flag.String("basic-auth.user", "", "")
var basicAuthPassword = flag.String("basic-auth.password", "", "")

func main() {
	flag.Parse()

	mrtgURL, err := url.Parse(*mrtgURL)
	if err != nil {
		panic(err)
	}
	auth := &mrtg.AuthenticateInfo{}
	auth.BasicAuth.User = *basicAuthUser
	auth.BasicAuth.Password = *basicAuthPassword

	result, err := mrtg.Scrape(context.Background(), mrtgURL, *timeout, auth)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
