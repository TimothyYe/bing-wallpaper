package bing_wallpaper

type Response struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	URL           string `json:"url"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyright_link"`
}
