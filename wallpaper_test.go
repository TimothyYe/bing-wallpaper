package bing_wallpaper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	resp, err := Get(context.Background(), 0, "zh-CN", "3840")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.URL)
}

func TestGetUnsupportedResolution(t *testing.T) {
	_, err := Get(context.Background(), 0, "zh-CN", "9999")
	assert.Error(t, err)
}
