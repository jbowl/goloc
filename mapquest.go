package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// LatLngResp  address and corresponding lattitude longitude
type LatLngResp struct {
	Address string  `json:"address"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
}

type MqAPI struct {
	Consumerkey string
}

//https://developer.mapquest.com/documentation/open/geocoding-api/address/get/
// following structs are to used parse the mapquest response
type Latlng struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type Locations struct {
	Type string `json:"type"`
	Ll   Latlng `json:"latLng"`
}

type ProvLocation struct {
	Location string `json:"location"`
}

type Results struct {
	ProvLoc ProvLocation `json:"providedLocation"`
	Locs    []Locations  `json:"locations"`
}

type ResultWrapper struct {
	Res []Results `json:"results"`
}

func mqRequest(url string) (llresp *LatLngResp, err error) {
	var client = &http.Client{
		Timeout: time.Second * 120,
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	//req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var rw ResultWrapper

	json.NewDecoder(resp.Body).Decode(&rw)
	if err != nil {
		return nil, err
	}

	llResp := LatLngResp{
		rw.Res[0].ProvLoc.Location,
		rw.Res[0].Locs[0].Ll.Lat,
		rw.Res[0].Locs[0].Ll.Lng}

	return &llResp, nil
}

// LatLng return a json struct holding geoloc lat and lng for location
//    use mapquest GEOLOC REST api to get lat lng
func (mq *MqAPI) LatLng(location string) (latlng *LatLngResp, err error) {

	baseURL, err := url.Parse("http://open.mapquestapi.com") ///geocoding/v1/address")
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return
	}
	baseURL.Path += "/geocoding/v1/address"

	params := url.Values{}
	params.Add("key", mq.Consumerkey)
	params.Add("location", location)
	baseURL.RawQuery = params.Encode() // Escape Query Parameters

	return mqRequest(baseURL.String())
}
