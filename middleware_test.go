package goapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriority(t *testing.T) {
	var m middlewares
	m = append(m, &middleware{priority: 1}, &middleware{priority: 4}, &middleware{priority: 0}, &middleware{priority: 10})
	sortByPriority(m)

	assert.Equal(t, 0, m[0].priority)
	assert.Equal(t, 1, m[1].priority)
	assert.Equal(t, 4, m[2].priority)
	assert.Equal(t, 10, m[3].priority)
}
