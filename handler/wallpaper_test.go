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
