package platforms

import "log"

const (
	vk = iota
)

var platformsMap = map[string]uint8{
	"VK": vk,
}

// GetByName return platform id
func GetByName(name string) uint8 {
	platformID, exist := platformsMap[name]
	if !exist {
		log.Fatalf("unknowwn platform %s", name)
	}

	return platformID
}
