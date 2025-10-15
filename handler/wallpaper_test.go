package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomIndex(t *testing.T) {
	randomIndex := getRandomIndex()
	t.Log("random index:", randomIndex)
	assert.GreaterOrEqual(t, randomIndex, 0)
	assert.LessOrEqual(t, randomIndex, 7)
}

func TestRandomMarket(t *testing.T) {
	validMarkets := map[string]bool{
		"zh-CN": true, "en-US": true, "ja-JP": true, "en-AU": true,
		"en-GB": true, "de-DE": true, "en-NZ": true, "en-CA": true,
		"en-IN": true, "fr-FR": true, "fr-CA": true, "it-IT": true,
		"es-ES": true, "pt-BR": true, "en-ROW": true,
	}

	randomMarket := getRandomMarket()
	t.Log("random market:", randomMarket)
	assert.True(t, validMarkets[randomMarket], "Expected valid market code, got: %s", randomMarket)
}
