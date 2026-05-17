package bing_wallpaper

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	bingURL   = `https://www.bing.com`
	bingAPI   = `https://www.bing.com/HPImageArchive.aspx?format=xml&idx=%d&n=1&mkt=%s`
	userAgent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0`

	// index 0 is "today" and can rotate at any time; older indices are
	// immutable so they can be cached aggressively.
	ttlCurrent  = 1 * time.Hour
	ttlArchived = 7 * 24 * time.Hour
)

var (
	Resolution     map[string]string
	FullResolution map[string]string

	httpClient = &http.Client{Timeout: 5 * time.Second}

	cacheMu sync.RWMutex
	cache   = make(map[string]cacheEntry)

	sfGroup singleflight.Group
)

type cacheEntry struct {
	resp    Response
	expires time.Time
}

type bingImage struct {
	StartDate     string `xml:"startdate"`
	EndDate       string `xml:"enddate"`
	URLBase       string `xml:"urlBase"`
	Copyright     string `xml:"copyright"`
	CopyrightLink string `xml:"copyrightlink"`
}

type bingResponse struct {
	XMLName xml.Name  `xml:"images"`
	Image   bingImage `xml:"image"`
}

func init() {
	Resolution = map[string]string{
		"1366": "1366x768.jpg",
		"1920": "1920x1080.jpg",
		"3840": "UHD.jpg",
	}

	FullResolution = map[string]string{
		"UHD":       "UHD.jpg",
		"1920x1200": "1920x1200.jpg",
		"1920x1080": "1920x1080.jpg",
		"1366x768":  "1366x768.jpg",
		"1280x768":  "1280x768.jpg",
		"1024x768":  "1024x768.jpg",
		"800x600":   "800x600.jpg",
		"800x480":   "800x480.jpg",
		"1080x1920": "1080x1920.jpg",
		"768x1280":  "768x1280.jpg",
		"720x1280":  "720x1280.jpg",
		"640x480":   "640x480.jpg",
		"480x800":   "480x800.jpg",
		"400x240":   "400x240.jpg",
		"320x240":   "320x240.jpg",
		"240x320":   "240x320.jpg",
	}
}

// cacheGet returns a value copy so callers can mutate without affecting
// the cached entry or other concurrent callers.
func cacheGet(key string) (Response, bool) {
	cacheMu.RLock()
	entry, ok := cache[key]
	cacheMu.RUnlock()
	if !ok || time.Now().After(entry.expires) {
		return Response{}, false
	}
	return entry.resp, true
}

func cacheSet(key string, resp Response, ttl time.Duration) {
	cacheMu.Lock()
	cache[key] = cacheEntry{resp: resp, expires: time.Now().Add(ttl)}
	cacheMu.Unlock()
}

// Get bing.com wallpaper from bing api
func Get(index uint, market, resolution string) (*Response, error) {
	suffix, ok := Resolution[resolution]
	if !ok {
		suffix, ok = FullResolution[resolution]
		if !ok {
			return nil, fmt.Errorf("resolution %s is not supported", resolution)
		}
	}
	resolution = suffix

	key := fmt.Sprintf("%d_%s_%s", index, market, resolution)

	if resp, ok := cacheGet(key); ok {
		return &resp, nil
	}

	v, err, _ := sfGroup.Do(key, func() (any, error) {
		// recheck after singleflight serialization in case another caller
		// already populated the cache while we waited.
		if resp, ok := cacheGet(key); ok {
			return resp, nil
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(bingAPI, index, market), nil)
		if err != nil {
			return Response{}, err
		}
		req.Header.Set("Referer", bingURL)
		req.Header.Set("User-Agent", userAgent)

		httpResp, err := httpClient.Do(req)
		if err != nil {
			return Response{}, err
		}
		defer httpResp.Body.Close()

		var br bingResponse
		if err := xml.NewDecoder(httpResp.Body).Decode(&br); err != nil {
			return Response{}, fmt.Errorf("failed to parse bing response: %w", err)
		}
		if br.Image.URLBase == "" {
			return Response{}, fmt.Errorf("empty image element in bing response")
		}

		response := Response{
			StartDate:     br.Image.StartDate,
			EndDate:       br.Image.EndDate,
			URL:           fmt.Sprintf("%s%s_%s", bingURL, br.Image.URLBase, resolution),
			Copyright:     br.Image.Copyright,
			CopyrightLink: br.Image.CopyrightLink,
		}

		ttl := ttlArchived
		if index == 0 {
			ttl = ttlCurrent
		}
		cacheSet(key, response, ttl)
		return response, nil
	})
	if err != nil {
		return nil, err
	}
	// singleflight shares the returned value across all callers; copy so
	// each caller gets its own pointer and downstream mutations cannot race.
	out := v.(Response)
	return &out, nil
}
