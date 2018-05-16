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

	debug = kingpin.Flag("debug", "log requests").Short('d').Counter()

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

	if *debug > 0 {
		ch.g.Use(apiserv.LogRequests(*debug > 1))
	}

	ch.g.GET("/userID", ch.GetUserID)
	ch.g.POST("/upgradeCampaign/:uid/:draftCID", ch.UpgradeCampaign)
	ch.g.GET("/adsReport/:uid/:start/:end", ch.GetAdsReport)
	ch.g.GET("/campaignReport/:uid/:cid/:start/:end", ch.GetCampaignReport)

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
			log.Printf("%T: %#+v", x, x)
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

func (ch *clientHandler) GetUserID(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.GetUserID(context.Background())
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) GetCampaignReport(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}
	var (
		uid, cid, start, end = ctx.Param("uid"), ctx.Param("cid"), ctx.Param("start"), ctx.Param("end")

		data interface{}
		err  error

		cacheKey, resp = ch.checkReportCache(ctx, c, uid, cid, start, end)
	)

	if resp != nil {
		return resp
	}

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
		uid, start, end = ctx.Param("uid"), ctx.Param("start"), ctx.Param("end")

		data interface{}
		err  error

		cacheKey, resp = ch.checkReportCache(ctx, c, uid, start, end)
	)

	if resp != nil {
		return resp
	}

	ch.c.Update(cacheKey, func(old interface{}) (_ interface{}, _ bool, _ time.Duration) {
		if data = old; data == nil {
			data, err = c.GetAdsReport(context.Background(), uid, sdk.DateToTime(start), sdk.DateToTime(end))
		}

		return data, false, time.Hour * 3
	})

	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) checkReportCache(ctx *apiserv.Context, c *sdk.Client, cacheKeyParts ...interface{}) (cacheKey string, resp apiserv.Response) {
	var (
		callsKey = fmt.Sprintf("client:%p", c)
		calls    uint64
	)

	ch.c.Update(callsKey, func(old interface{}) (_ interface{}, _ bool, _ time.Duration) {
		calls, _ := old.(uint64)
		calls++
		return calls, true, time.Hour
	})

	if calls > maxReportingCallsPerHour {
		resp = apiserv.NewJSONErrorResponse(http.StatusTooManyRequests,
			fmt.Sprintf("you went over your hourly request limit of %d by %d",
				maxReportingCallsPerHour, calls-maxReportingCallsPerHour,
			))
		return
	}

	cacheKey = cache.Key(cacheKeyParts...)
	if data, ok := ch.c.Get(cacheKey); ok {
		resp = apiserv.NewJSONResponse(data)
		return
	}

	return
}

func (ch *clientHandler) UpgradeCampaign(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.UpgradeCampaign(context.Background(), ctx.Param("uid"), ctx.Param("draftCID"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}
