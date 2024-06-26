package html

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReportHtmlTemplatesFunctions(t *testing.T) {
	t.Run("test increment function", func(t *testing.T) {
		number := 1234

		result := increment()(number)

		assert.Equal(t, 1235, result)
	})

	t.Run("test checkStatus function", func(t *testing.T) {
		result := checkStatus()("PASS")
		assert.True(t, result)

		result = checkStatus()("YES")
		assert.True(t, result)

		result = checkStatus()("foobar")
		assert.False(t, result)
	})

	t.Run("test calculateRatio function", func(t *testing.T) {
		result := calculateRatio()(10, 100)
		assert.Equal(t, 10, result)

		result = calculateRatio()(0, 45)
		assert.Equal(t, 0, result)

		result = calculateRatio()(57846, 0)
		assert.Equal(t, 100, result)
	})

	t.Run("test formatDateTime function", func(t *testing.T) {
		inputTime := time.Time{}.AddDate(2000, 4, 21).Add(time.Hour * 13).Add(time.Minute * 24).Add(time.Second * 49)
		result := formatDateTime()(inputTime)
		assert.Equal(t, "2001/05/22 13:24:49", result)
	})

	t.Run("test formatDate function", func(t *testing.T) {
		inputTime := time.Time{}.AddDate(2000, 4, 21).Add(time.Hour * 13).Add(time.Minute * 24).Add(time.Second * 49)
		result := formatDate()(inputTime)
		assert.Equal(t, "2001/05/22", result)
	})

	t.Run("test formatTime function", func(t *testing.T) {
		inputTime := time.Time{}.AddDate(2000, 4, 21).Add(time.Hour * 13).Add(time.Minute * 24).Add(time.Second * 49)
		result := formatTime()(inputTime)
		assert.Equal(t, "13:24:49", result)
	})

	t.Run("test toHumanTime function", func(t *testing.T) {
		result := toHumanTime()(time.Second * 23)
		assert.Equal(t, "23 [s]", result)

		result = toHumanTime()(time.Millisecond * 23)
		assert.Equal(t, "23 [ms]", result)

		result = toHumanTime()(time.Microsecond * 23)
		assert.Equal(t, "23 [μs]", result)

		result = toHumanTime()(time.Nanosecond * 23)
		assert.Equal(t, "23 [ns]", result)

		result = toHumanTime()(time.Hour * 40)
		assert.Equal(t, "144000 [s]", result)
	})
}
