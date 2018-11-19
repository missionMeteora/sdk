package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/PathDNA/ptk"
	"github.com/PathDNA/ptk/cache"
	"github.com/PathDNA/ssync"
	"github.com/missionMeteora/apiserv"
	"github.com/missionMeteora/sdk"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "sdk " + sdk.Version + " (server v0.3)"

	clientCacheTimeout       = time.Hour * 24
	maxReportingCallsPerHour = 100
)

var (
	live    = kingpin.Flag("live", "are we using the live api endpoints").Short('l').Bool()
	apiAddr = kingpin.Flag("apiAddr", "local api addr").Default("http://localhost:8080").Short('a').String()

	debug = kingpin.Flag("debug", "log requests").Short('d').Counter()

	addr       = kingpin.Flag("addr", "listen addr").Default(":8081").String()
	ssyncAddr  = kingpin.Flag("ssyncAddr", "ssync addr").String()
	clicksAddr = kingpin.Flag("clicksAddr", "clicks addr").String()
	visitsAddr = kingpin.Flag("visitsAddr", "visits addr").String()
	pathAddr   = kingpin.Flag("pathAddr", "path addr").String()

	letsEnc = kingpin.Flag("letsencrypt", "run production letsencrypt, addr must be set to a valid hostname").Short('s').Bool()

	apiPrefix = kingpin.Flag("prefix", "api route prefix").Default("/api/v1").Short('p').String()

	pongResp = apiserv.NewJSONResponse("pong")
	verResp  = apiserv.NewJSONResponse(version)

	pathIDsMap atomic.Value
)

func main() {
	log.SetFlags(log.Lshortfile)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version).VersionFlag.Short('V')
	kingpin.Parse()

	sc, err := ssync.NewClient("", "bank", *ssyncAddr)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserv.New(apiserv.SetNoCatchPanics(true))
	ch := &clientHandler{
		g:  s.Group(*apiPrefix),
		sc: sc,
		c:  cache.NewMemCache(time.Minute * 15),
	}

	if *debug > 0 {
		ch.g.Use(apiserv.LogRequests(*debug > 1))
	}

	ch.g.GET("/userID", ch.GetUserID)
	ch.g.GET("/listApps", ch.listApps)
	ch.g.POST("/upgradeCampaign/:uid/:draftCID", ch.UpgradeCampaign)
	ch.g.GET("/adsReport/:uid/:start/:end", ch.GetAdsReport)
	ch.g.GET("/campaignReport/:uid/:cid/:start/:end", ch.GetCampaignReport)
	ch.g.GET("/receipts/:uid/:date", ch.GetReceipts)
	ch.g.GET("/receipts/:uid/:cid/:date", ch.GetReceipts)
	ch.g.GET("/clicks/:uid/:date", ch.GetClicks)
	ch.g.GET("/clicks/:uid/:cid/:date", ch.GetClicks)
	ch.g.GET("/visits/:uid/:date", ch.GetVisits)
	ch.g.GET("/visits/:uid/:cid/:date", ch.GetVisits)

	ch.g.GET("/ping", func(*apiserv.Context) apiserv.Response { return pongResp })
	ch.g.GET("/version", func(*apiserv.Context) apiserv.Response { return verResp })

	if *live {
		ch.addr = sdk.DefaultServer
	} else {
		ch.addr = *apiAddr
	}

	if *pathAddr != "" {
		log.Printf("running path ids mapper")
		go updatePathIDsMap(*pathAddr)
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

	sc *ssync.Client
	c  *cache.MemCache
	g  apiserv.Group
}

func (ch *clientHandler) getClient(ctx *apiserv.Context) (c *sdk.Client) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%T: %v", x, x)
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

func (ch *clientHandler) listApps(ctx *apiserv.Context) apiserv.Response {
	ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	type appWithName struct {
		Name string  `json:"name"`
		App  sdk.App `json:"app"`
	}

	var allApps []appWithName
	for _, app := range sdk.AllApps {
		allApps = append(allApps, appWithName{app.Name(), app})
	}

	ctx.JSON(200, true, allApps)
	return nil
}

func (ch *clientHandler) GetReceipts(ctx *apiserv.Context) apiserv.Response {
	var (
		c = ch.getClient(ctx)

		uid  = ctx.Param("uid")
		cid  = ctx.Param("cid")
		date = sdk.DateToTime(ctx.Param("date"))
	)

	receipts, err := c.Receipts(ctx.Req.Context(), ch.sc, date, uid, cid)
	if err != nil {
		log.Printf("receipts %s/%s/%s:%v", uid, cid, ctx.Param("date"), err)
		return apiserv.NewJSONErrorResponse(500)
	}

	return apiserv.PlainResponse(apiserv.MimeJSON, receipts)
}

func (ch *clientHandler) GetClicks(ctx *apiserv.Context) apiserv.Response {
	var (
		c = ch.getClient(ctx)

		uid  = ctx.Param("uid")
		cid  = ctx.Param("cid")
		date = sdk.DateToTime(ctx.Param("date"))
	)

	clicks, err := c.Clicks(ctx.Req.Context(), *clicksAddr, date, uid, cid)
	if err != nil {
		log.Printf("clicks %s/%s/%s:%v", uid, cid, ctx.Param("date"), err)
		return apiserv.NewJSONErrorResponse(500)
	}

	return apiserv.PlainResponse(apiserv.MimeJSON, clicks)
}

func (ch *clientHandler) GetVisits(ctx *apiserv.Context) apiserv.Response {
	var (
		c = ch.getClient(ctx)

		uid  = ctx.Param("uid")
		cid  = ctx.Param("cid")
		date = sdk.DateToTime(ctx.Param("date"))
	)

	visits, err := c.Visits(ctx.Req.Context(), *visitsAddr, date, uid, cid)
	if err != nil {
		log.Printf("visits %s/%s/%s:%v", uid, cid, ctx.Param("date"), err)
		return apiserv.NewJSONErrorResponse(500)
	}

	return apiserv.PlainResponse(apiserv.MimeJSON, visits)
}

func updatePathIDsMap(url string) {
	for {
		var resp map[string]string
		ptk.Request("GET", "", url, nil, &resp)
		if len(resp) > 0 {
			m := make(map[string]string, len(resp))
			for pid, mid := range resp {
				m[mid] = min(m[mid], pid)
			}
			pathIDsMap.Store(m)
		}
		time.Sleep(time.Minute * 30)
	}
}

func min(a, b string) string {
	ai, bi := atoi(a), atoi(b)
	if ai > 0 && ai < bi {
		return a
	}
	return b
}

func atoi(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}
