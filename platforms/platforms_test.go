package platforms_test

import (
	"testing"

	"github.com/server-may-cry/bubble-go/platforms"
)

var platformsTests = []struct {
	platform string // input
	expected uint8  // expected result
}{
	{"VK", 1},
	{"OK", 2},
}

func TestGetByName(t *testing.T) {
	for _, tt := range platformsTests {
		actual := platforms.GetByName(tt.platform)
		if actual != tt.expected {
			t.Errorf("platforms.GetByName(%s): expected %d, actual %d", tt.platform, tt.expected, actual)
		}
	}
}
