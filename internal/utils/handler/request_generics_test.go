package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.strale.io/go-travel/internal/utils/handler"
)

func TestMyCode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/temp?id=123", nil)
	i, err := handler.Query(req, handler.Int64, "id", true, 1, handler.IsPositive)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), i.Val())
}
