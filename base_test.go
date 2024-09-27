package base_test

import (
	"io"
	"os"
	"testing"

	"github.com/baldanca/go-base"
	"github.com/stretchr/testify/assert"
)

type envType struct{}

func TestNewBaseDefaultConfig(t *testing.T) {
	assert.Equal(t, "{\"level\":\"info\",\"message\":\"stdout test\"}\n", captureOutput(func() {
		newBase := base.New[envType](base.Config{})
		assert.NotNil(t, newBase)

		assert.NotNil(t, newBase.Ctx())
		assert.NotNil(t, newBase.CancelCtx())

		assert.NotNil(t, newBase.Env())

		assert.NotNil(t, newBase.HTTPClient())
		assert.Equal(t, newBase.HTTPClient().Timeout, base.DefaultHTTPClientTimeout)

		assert.NotNil(t, newBase.Logger())
		assert.Equal(t, newBase.Logger().GetLevel(), base.DefaultLoggerLevel)
		newBase.Logger().Info().Msg("stdout test")

		assert.NotNil(t, newBase.TimeLocation())
		assert.Equal(t, newBase.TimeLocation().String(), base.DefaultTimeLocation)

		assert.NotNil(t, newBase.TimeNow())
		assert.Greater(t, newBase.TimeNow().Unix(), int64(0))
		assert.Equal(t, newBase.TimeNow().Location(), newBase.TimeLocation())
	}))
}

func captureOutput(f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}
