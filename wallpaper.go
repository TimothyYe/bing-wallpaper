package bing_wallpaper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/beevik/etree"
)

const (
	bingURL = `https://www.bing.com`
	bingAPI = `https://www.bing.com/HPImageArchive.aspx?format=xml&idx=%d&n=1&mkt=%s`
)

var (
	Resolution map[string]string
	Markets    map[string]bool
	cache      *bigcache.BigCache
)

func init() {
	Resolution = map[string]string{
		"1366": "1366x768.jpg",
		"1920": "1920x1080.jpg",
		"3840": "UHD.jpg",
	}

	Markets = map[string]bool{
		"en-US": true,
		"zh-CN": true,
		"ja-JP": true,
		"en-AU": true,
		"en-GB": true,
		"de-DE": true,
		"en-NZ": true,
		"en-CA": true,
		"en-IN": true,
		"fr-FR": true,
		"fr-CA": true,
	}

	// initialize the cache
	config := bigcache.Config{
		Shards:             128,
		LifeWindow:         60 * time.Minute,
		CleanWindow:        10 * time.Minute,
		MaxEntriesInWindow: 30 * 60,
		MaxEntrySize:       50,
		Verbose:            true,
		HardMaxCacheSize:   256,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}

	var initErr error
	cache, initErr = bigcache.New(context.Background(), config)
	if initErr != nil {
		log.Fatal(initErr)
	}
}

// Get bing.com wallpaper from bing api
func Get(index uint, market, resolution string) (*Response, error) {
	if _, ok := Resolution[resolution]; !ok {
		return nil, fmt.Errorf("resolution %s is not supported", resolution)
	}

	if _, ok := Markets[market]; !ok {
		return nil, fmt.Errorf("market %s is not supported", market)
	}

	// query cache first
	if value, err := cache.Get(fmt.Sprintf("%d_%s_%s", index, market, resolution)); err == nil {
		cachedResp := &Response{}
		_ = json.Unmarshal(value, cachedResp)
		return cachedResp, nil
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(bingAPI, index, market), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Referer", bingURL)
	request.Header.Add("User-Agent", `Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 10.0; WOW64; Trident/8.0; .NET4.0C; .NET4.0E)`)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request body from %s", bingURL)
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		return nil, err
	}

	// get image element
	imgElem := doc.SelectElement("images").SelectElement("image")

	response := &Response{
		StartDate:     imgElem.SelectElement("startdate").Text(),
		EndDate:       imgElem.SelectElement("enddate").Text(),
		URL:           fmt.Sprintf("%s%s_%s", bingURL, imgElem.SelectElement("urlBase").Text(), Resolution[resolution]),
		Copyright:     imgElem.SelectElement("copyright").Text(),
		CopyrightLink: imgElem.SelectElement("copyrightlink").Text(),
	}

	// cache the response
	if value, err := json.Marshal(response); err == nil {
		_ = cache.Set(fmt.Sprintf("%d_%s_%s", index, market, resolution), value)
	}

	return response, nil
}
