package evaluate

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetRandomStampResponse(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run("random", func(t *testing.T) {
			got, err := GetRandomStampResponse()
			assert.Equal(t, err, nil)
			fmt.Println(got)
			assert.Equal(t, got != "", true)
		})
	}
}
