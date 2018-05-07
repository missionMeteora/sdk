package main

import (
	"net/url"
	"log"
	"net/http"
	"sync"

	"github.com/missionMeteora/apiserv"
	"github.com/missionMeteora/sdk"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "sdk " + sdk.Version + ", server v0.1"

var (
	live = kingpin.Flag("live", "are we using the live api endpoints").Short('l').Bool()
	apiAddr = kingpin.Flag("apiAddr", "local api addr").Default("http://localhost:8080").Short('a').String()

	addr = kingpin.Flag("addr", "listen addr").Default(":8081").String()

	letsEnc = kingpin.Flag("letsencrypt", "run production letsencrypt, addr must be set to a valid hostname").Short('s').Bool()

	apiPrefix = kingpin.Flag("prefix", "api route prefix").Default("/api/v1").Short('p').String()

	pongResp = apiserv.NewJSONResponse("pong")
	verResp = apiserv.NewJSONResponse(version)
)

func main() {
	log.SetFlags(log.Lshortfile)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version).VersionFlag.Short('V')
	kingpin.Parse()

	s := apiserv.New()

	ch := &clientHandler{
		g: s.Group(*apiPrefix),
	}

	ch.g.GET("/ping", func(*apiserv.Context) apiserv.Response { return pongResp })
	ch.g.GET("/version", func(*apiserv.Context) apiserv.Response { return verResp })

	if *live {
		ch.addr = sdk.DefaultServer
	} else {
		ch.addr = *apiAddr
	}
	ch.init()

	if *letsEnc {
		u, err := url.Parse(*addr)
		if err != nil {
			log.Fatalf("invalid host: %s", *addr)
		}

		host := u.Hostname()
		log.Printf("Listening on https://%s", host)
		log.Fatal(s.RunAutoCert("", host))
	} else {
		log.Printf("Listening on http://%s", *addr)
		log.Fatal(s.Run(*addr))
	}
}

type clientHandler struct{
	addr string
	g apiserv.Group
	m sync.Map
}

func (ch *clientHandler) getClient(ctx *apiserv.Context) (_ *sdk.Client) {
	defer func() {
		if x := recover(); x != nil {
			ctx.JSON(http.StatusInternalServerError, true, x)
		}
	}()

	key := ctx.Query("apiKey")
	if key == "" {
		key = ctx.Header().Get("X-APIKEY")
	}

	if key == "" {
		ctx.JSON(http.StatusUnauthorized, true, "missing api key")
		return
	}

	if c, ok := ch.m.Load(key); ok {
		return c.(*sdk.Client)
	}

	c, _ := ch.m.LoadOrStore(key, sdk.NewWithAddr(ch.addr, key))
	return c.(*sdk.Client)
}
