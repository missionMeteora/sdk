package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PathDNA/ptk/cache"
	"github.com/missionMeteora/apiserv"
	"github.com/missionMeteora/sdk"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "sdk " + sdk.Version + " (server v0.1)"

	clientCacheTimeout       = time.Hour * 24
	maxReportingCallsPerHour = 100
)

var (
	live    = kingpin.Flag("live", "are we using the live api endpoints").Short('l').Bool()
	apiAddr = kingpin.Flag("apiAddr", "local api addr").Default("http://localhost:8080").Short('a').String()

	addr = kingpin.Flag("addr", "listen addr").Default(":8081").String()

	letsEnc = kingpin.Flag("letsencrypt", "run production letsencrypt, addr must be set to a valid hostname").Short('s').Bool()

	apiPrefix = kingpin.Flag("prefix", "api route prefix").Default("/api/v1").Short('p').String()

	pongResp = apiserv.NewJSONResponse("pong")
	verResp  = apiserv.NewJSONResponse(version)
)

func main() {
	log.SetFlags(log.Lshortfile)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version).VersionFlag.Short('V')
	kingpin.Parse()

	s := apiserv.New()

	ch := &clientHandler{
		g: s.Group(*apiPrefix),
		c: cache.NewMemCache(time.Minute * 15),
	}

	ch.g.AddRoute("GET", "/adsReport/:uid/:start/:end", ch.GetAdsReport)
	ch.g.AddRoute("GET", "/campaignReport/:uid/:cid/:start/:end", ch.GetCampaignReport)

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

type clientHandler struct {
	addr string

	c *cache.MemCache
	g apiserv.Group
}

func (ch *clientHandler) getClient(ctx *apiserv.Context) (c *sdk.Client) {
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
		ctx.JSON(http.StatusUnauthorized, true, apiserv.NewJSONErrorResponse(http.StatusUnauthorized, "missing or invalid api key"))
		return
	}

	if c, ok := ch.c.Get(key); ok {
		return c.(*sdk.Client)
	}

	var err error
	ch.c.Update(key, func(old interface{}) (_ interface{}, _ bool, _ time.Duration) {
		if c, _ = old.(*sdk.Client); c != nil {
			return c, false, time.Hour * 24
		}

		c = sdk.NewWithAddr(ch.addr, key)
		var uid string
		if uid, err = c.GetUserID(context.Background()); err != nil {
			return
		}

		// ghetto check, this needs to be handled by an actual endpoint
		if _, err = c.ListAdvertisers(context.Background(), uid); err != nil {
			return
		}

		return c, false, clientCacheTimeout
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, true, apiserv.NewJSONErrorResponse(http.StatusUnauthorized, err))
		c = nil
	}
	return
}

func (ch *clientHandler) GetCampaignReport(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}
	var (
		callsKey             = fmt.Sprintf("client:%p", c)
		uid, cid, start, end = ctx.Param("uid"), ctx.Param("cid"), ctx.Param("start"), ctx.Param("end")
		calls                uint64
	)

	ch.c.Update(callsKey, func(old interface{}) (_ interface{}, _ bool, _ time.Duration) {
		calls, _ := old.(uint64)
		calls++
		return calls, true, time.Hour
	})

	if calls > maxReportingCallsPerHour {
		return apiserv.NewJSONErrorResponse(http.StatusTooManyRequests,
			fmt.Sprintf("you went over your hourly request limit of %d by %d",
				maxReportingCallsPerHour, calls-maxReportingCallsPerHour,
			))
	}

	cacheKey := cache.Key(uid, cid, start, end)
	if data, ok := ch.c.Get(cacheKey); ok {
		return apiserv.NewJSONResponse(data)
	}

	var (
		data interface{}
		err  error
	)

	ch.c.Update(cacheKey, func(old interface{}) (_ interface{}, _ bool, _ time.Duration) {
		if data = old; data == nil {
			data, err = c.GetCampaignReport(context.Background(), uid, cid, sdk.DateToTime(start), sdk.DateToTime(end))
		}

		return data, false, time.Hour * 3
	})

	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) GetAdsReport(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var (
		callsKey        = fmt.Sprintf("client:%p", c)
		uid, start, end = ctx.Param("uid"), ctx.Param("start"), ctx.Param("end")
		calls           uint64
	)

	ch.c.Update(callsKey, func(old interface{}) (_ interface{}, _ bool, _ time.Duration) {
		calls, _ := old.(uint64)
		calls++
		return calls, true, time.Hour
	})

	if calls > maxReportingCallsPerHour {
		return apiserv.NewJSONErrorResponse(http.StatusTooManyRequests,
			fmt.Sprintf("you went over your hourly request limit of %d by %d",
				maxReportingCallsPerHour, calls-maxReportingCallsPerHour,
			))
	}

	cacheKey := cache.Key(uid, start, end)
	if data, ok := ch.c.Get(cacheKey); ok {
		return apiserv.NewJSONResponse(data)
	}

	var (
		data interface{}
		err  error
	)

	log.Println(uid, sdk.DateToTime(start), sdk.DateToTime(end))

	ch.c.Update(cacheKey, func(old interface{}) (_ interface{}, _ bool, _ time.Duration) {
		if data = old; data == nil {
			data, err = c.GetAdsReport(context.Background(), uid, sdk.DateToTime(start), sdk.DateToTime(end))
			log.Println(err)
		}

		return data, false, time.Hour * 3
	})

	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}
