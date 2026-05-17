package bing_wallpaper

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
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

	// ttlNegative briefly caches upstream failures so a Bing outage or
	// 5xx burst does not get amplified by every incoming request.
	ttlNegative = 5 * time.Second
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
	err     error // non-nil means this is a negative-cache entry
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
// the cached entry or other concurrent callers. If the entry is a
// negative-cache entry, err is non-nil and resp is the zero value.
func cacheGet(key string) (Response, error, bool) {
	cacheMu.RLock()
	entry, ok := cache[key]
	cacheMu.RUnlock()
	if !ok || time.Now().After(entry.expires) {
		return Response{}, nil, false
	}
	return entry.resp, entry.err, true
}

func cacheSet(key string, resp Response, err error, ttl time.Duration) {
	cacheMu.Lock()
	cache[key] = cacheEntry{resp: resp, err: err, expires: time.Now().Add(ttl)}
	cacheMu.Unlock()
}

// Get bing.com wallpaper from bing api. The provided context is used for
// the upstream HTTP call so a client disconnect cancels the in-flight
// request to Bing.
func Get(ctx context.Context, index uint, market, resolution string) (*Response, error) {
	suffix, ok := Resolution[resolution]
	if !ok {
		suffix, ok = FullResolution[resolution]
		if !ok {
			return nil, fmt.Errorf("resolution %s is not supported", resolution)
		}
	}
	resolution = suffix

	key := fmt.Sprintf("%d_%s_%s", index, market, resolution)

	if resp, cachedErr, ok := cacheGet(key); ok {
		if cachedErr != nil {
			return nil, cachedErr
		}
		return &resp, nil
	}

	v, err, _ := sfGroup.Do(key, func() (any, error) {
		// recheck after singleflight serialization in case another caller
		// already populated the cache while we waited.
		if resp, cachedErr, ok := cacheGet(key); ok {
			if cachedErr != nil {
				return Response{}, cachedErr
			}
			return resp, nil
		}

		response, fetchErr := fetchFromBing(ctx, index, market, resolution)
		if fetchErr != nil {
			cacheSet(key, Response{}, fetchErr, ttlNegative)
			return Response{}, fetchErr
		}

		ttl := ttlArchived
		if index == 0 {
			ttl = ttlCurrent
		}
		cacheSet(key, response, nil, ttl)
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

func fetchFromBing(ctx context.Context, index uint, market, resolution string) (Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(bingAPI, index, market), nil)
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Referer", bingURL)
	req.Header.Set("User-Agent", userAgent)

	httpResp, err := httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		// drain so the connection can be reused from the pool
		_, _ = io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return Response{}, fmt.Errorf("bing returned HTTP %d", httpResp.StatusCode)
	}

	var br bingResponse
	if err := xml.NewDecoder(httpResp.Body).Decode(&br); err != nil {
		return Response{}, fmt.Errorf("failed to parse bing response: %w", err)
	}
	if br.Image.URLBase == "" {
		return Response{}, fmt.Errorf("empty image element in bing response")
	}

	return Response{
		StartDate:     br.Image.StartDate,
		EndDate:       br.Image.EndDate,
		URL:           fmt.Sprintf("%s%s_%s", bingURL, br.Image.URLBase, resolution),
		Copyright:     br.Image.Copyright,
		CopyrightLink: br.Image.CopyrightLink,
	}, nil
}
