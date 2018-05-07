// this file is automatically generated, make sure you don't ovewrite your changes

package main

import (
	"context"
	"net/http"

	"github.com/missionMeteora/apiserv"
	"github.com/missionMeteora/sdk"
)

/*
// add this to main.go and handle the logic in it
// it may return errors using ctx
// remember to call ch.init(g)
// these are notes to myself, don't judge, I barely remember my name, ok?
type clientHandler struct{}
func (ch *clientHandler) getClient(ctx *apiserv.Context) *sdk.Client {
		return nil
}
*/

func (ch *clientHandler) CreateAd(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var req *sdk.CreateAdRequest
	if err := ctx.BindJSON(&req); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := c.CreateAd(context.Background(), ctx.Param("uid"), req)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) CreateAdGroup(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.CreateAdGroup(context.Background(), ctx.Param("uid"), ctx.Param("name"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) CreateAdvertiser(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var req *sdk.CreateAdvertiserRequest
	if err := ctx.BindJSON(&req); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := c.CreateAdvertiser(context.Background(), req)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) CreateCampaign(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var cmp *sdk.Campaign
	if err := ctx.BindJSON(&cmp); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := c.CreateCampaign(context.Background(), ctx.Param("uid"), cmp)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) CreateDraftCampaign(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var cmp *sdk.Campaign
	if err := ctx.BindJSON(&cmp); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := c.CreateDraftCampaign(context.Background(), ctx.Param("uid"), cmp)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) CreateFullCampaign(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var req *sdk.CreateFullCampaignRequest
	if err := ctx.BindJSON(&req); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := c.CreateFullCampaign(context.Background(), ctx.Param("uid"), req)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) CreateProximitySegment(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var seg *sdk.ProximitySegment
	if err := ctx.BindJSON(&seg); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := c.CreateProximitySegment(context.Background(), ctx.Param("uid"), seg)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) CreateSegment(ctx *apiserv.Context) apiserv.Response { // method:POST
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var seg *sdk.Segment
	if err := ctx.BindJSON(&seg); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := c.CreateSegment(context.Background(), ctx.Param("uid"), seg)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) DeleteAd(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := "OK", c.DeleteAd(context.Background(), ctx.Param("adID"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) DeleteAdGroup(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := "OK", c.DeleteAdGroup(context.Background(), ctx.Param("adgroupID"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) DeleteCampaign(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := "OK", c.DeleteCampaign(context.Background(), ctx.Param("cid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) DeleteDraftCampaign(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := "OK", c.DeleteDraftCampaign(context.Background(), ctx.Param("cid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) DeleteProximitySegment(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := "OK", c.DeleteProximitySegment(context.Background(), ctx.Param("segID"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) DeleteSegment(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := "OK", c.DeleteSegment(context.Background(), ctx.Param("segID"))
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

	data, err := c.GetAdsReport(context.Background(), ctx.Param("uid"), ctx.Param("start"), ctx.Param("end"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) GetCampaign(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.GetCampaign(context.Background(), ctx.Param("cid"))
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

	data, err := c.GetCampaignReport(context.Background(), ctx.Param("uid"), ctx.Param("cid"), ctx.Param("start"), ctx.Param("end"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) GetDraftCampaign(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.GetDraftCampaign(context.Background(), ctx.Param("cid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) GetHeatmap(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.GetHeatmap(context.Background(), ctx.Param("uid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) GetUserAPIKey(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.GetUserAPIKey(context.Background(), ctx.Param("uid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListAdGroups(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListAdGroups(context.Background(), ctx.Param("uid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListAds(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListAds(context.Background(), ctx.Param("uid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListAdsByAdGroup(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListAdsByAdGroup(context.Background(), ctx.Param("uid"), ctx.Param("adGroupID"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListAdvertisers(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListAdvertisers(context.Background(), ctx.Param("agencyID"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListAgencies(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListAgencies(context.Background())
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListCampaigns(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListCampaigns(context.Background(), ctx.Param("uid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListDraftCampaigns(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListDraftCampaigns(context.Background(), ctx.Param("uid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListProximitySegments(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListProximitySegments(context.Background(), ctx.Param("uid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) ListSegments(ctx *apiserv.Context) apiserv.Response {
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	data, err := c.ListSegments(context.Background(), ctx.Param("uid"))
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) UpdateAd(ctx *apiserv.Context) apiserv.Response { // method:PUT
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var ad *sdk.Ad
	if err := ctx.BindJSON(&ad); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := "OK", c.UpdateAd(context.Background(), ad)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) UpdateCampaign(ctx *apiserv.Context) apiserv.Response { // method:PUT
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var cmp *sdk.Campaign
	if err := ctx.BindJSON(&cmp); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := "OK", c.UpdateCampaign(context.Background(), cmp)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) UpdateDraftCampaign(ctx *apiserv.Context) apiserv.Response { // method:PUT
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var cmp *sdk.Campaign
	if err := ctx.BindJSON(&cmp); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := "OK", c.UpdateDraftCampaign(context.Background(), cmp)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) UpdateProximitySegment(ctx *apiserv.Context) apiserv.Response { // method:PUT
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var seg *sdk.ProximitySegment
	if err := ctx.BindJSON(&seg); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := "OK", c.UpdateProximitySegment(context.Background(), seg)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) UpdateSegment(ctx *apiserv.Context) apiserv.Response { // method:PUT
	c := ch.getClient(ctx)
	if ctx.Done() {
		return nil
	}

	var seg *sdk.Segment
	if err := ctx.BindJSON(&seg); err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	data, err := "OK", c.UpdateSegment(context.Background(), seg)
	if err != nil {
		return apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)
	}

	return apiserv.NewJSONResponse(data)
}

func (ch *clientHandler) init() {
	ch.g.AddRoute("POST", "/ad", ch.CreateAd)
	ch.g.AddRoute("POST", "/adGroup", ch.CreateAdGroup)
	ch.g.AddRoute("POST", "/advertiser", ch.CreateAdvertiser)
	ch.g.AddRoute("POST", "/campaign", ch.CreateCampaign)
	ch.g.AddRoute("POST", "/draftCampaign", ch.CreateDraftCampaign)
	ch.g.AddRoute("POST", "/fullCampaign", ch.CreateFullCampaign)
	ch.g.AddRoute("POST", "/proximitySegment", ch.CreateProximitySegment)
	ch.g.AddRoute("POST", "/segment", ch.CreateSegment)
	ch.g.AddRoute("DELETE", "/ad", ch.DeleteAd)
	ch.g.AddRoute("DELETE", "/adGroup", ch.DeleteAdGroup)
	ch.g.AddRoute("DELETE", "/campaign", ch.DeleteCampaign)
	ch.g.AddRoute("DELETE", "/draftCampaign", ch.DeleteDraftCampaign)
	ch.g.AddRoute("DELETE", "/proximitySegment", ch.DeleteProximitySegment)
	ch.g.AddRoute("DELETE", "/segment", ch.DeleteSegment)
	ch.g.AddRoute("GET", "/adsReport", ch.GetAdsReport)
	ch.g.AddRoute("GET", "/campaign", ch.GetCampaign)
	ch.g.AddRoute("GET", "/campaignReport", ch.GetCampaignReport)
	ch.g.AddRoute("GET", "/draftCampaign", ch.GetDraftCampaign)
	ch.g.AddRoute("GET", "/heatmap", ch.GetHeatmap)
	ch.g.AddRoute("GET", "/userAPIKey", ch.GetUserAPIKey)
	ch.g.AddRoute("GET", "/adGroups", ch.ListAdGroups)
	ch.g.AddRoute("GET", "/ads", ch.ListAds)
	ch.g.AddRoute("GET", "/adsByAdGroup", ch.ListAdsByAdGroup)
	ch.g.AddRoute("GET", "/advertisers", ch.ListAdvertisers)
	ch.g.AddRoute("GET", "/agencies", ch.ListAgencies)
	ch.g.AddRoute("GET", "/campaigns", ch.ListCampaigns)
	ch.g.AddRoute("GET", "/draftCampaigns", ch.ListDraftCampaigns)
	ch.g.AddRoute("GET", "/proximitySegments", ch.ListProximitySegments)
	ch.g.AddRoute("GET", "/segments", ch.ListSegments)
	ch.g.AddRoute("PUT", "/ad", ch.UpdateAd)
	ch.g.AddRoute("PUT", "/campaign", ch.UpdateCampaign)
	ch.g.AddRoute("PUT", "/draftCampaign", ch.UpdateDraftCampaign)
	ch.g.AddRoute("PUT", "/proximitySegment", ch.UpdateProximitySegment)
	ch.g.AddRoute("PUT", "/segment", ch.UpdateSegment)
}
