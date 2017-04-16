package platforms

import (
	"testing"
)

var platformsTests = []struct {
	platform string // input
	expected uint8  // expected result
}{
	{"VK", 1},
	{"OK", 2},
}

func TestPlatformGetByName(t *testing.T) {
	for _, tt := range platformsTests {
		actual := GetByName(tt.platform)
		if actual != tt.expected {
			t.Errorf("platforms.GetByName(%s): expected %d, actual %d", tt.platform, tt.expected, actual)
		}
	}
}

func TestPlatformNotExist(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic expected on platform %s", "FB")
		}
	}()
	GetByName("FB")
}
