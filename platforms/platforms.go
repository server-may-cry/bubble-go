package platforms

import "fmt"

const (
	vk = iota + 1 // 1
	ok            // 2
)

var platformsMap = map[string]uint8{
	"VK": vk,
	"OK": ok,
}

// GetByName return platform id
func GetByName(name string) uint8 {
	platformID, exist := platformsMap[name]
	if !exist {
		panic(fmt.Sprintf("unknown platform %s", name))
	}

	return platformID
}
