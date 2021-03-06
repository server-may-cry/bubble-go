package platforms

import (
	"testing"
)

func TestPlatformGetByName(t *testing.T) {
	var platformsTests = [...]struct {
		platform string // input
		expected uint8  // expected result
		exist    bool
	}{
		{"VK", 1, true},
		{"OK", 2, true},
		{"FB", 0, false},
	}
	for _, tt := range platformsTests {
		actual, exist := GetByName(tt.platform)
		if exist != tt.exist {
			t.Errorf("platforms.GetByName(%s): exist expected %t, got %t", tt.platform, tt.exist, exist)
		}
		if actual != tt.expected {
			t.Errorf("platforms.GetByName(%s): expected %d, got %d", tt.platform, tt.expected, actual)
		}
	}
}
