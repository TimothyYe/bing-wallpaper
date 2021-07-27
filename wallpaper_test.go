package bing_wallpaper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	resp, err := Get(0, "zh-CN", "3840")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.URL)
}

func TestGetCauseErr(t *testing.T) {
	_, err := Get(0, "zh-TW", "1920")
	assert.Error(t, err)
}
