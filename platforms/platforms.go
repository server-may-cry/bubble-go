package platforms

const (
	// VK id in DB
	VK = 1
	// OK id in DB
	OK = 2
)

var platformsMap = map[string]uint8{
	"VK": VK,
	"OK": OK,
}

// GetByName return platform id
func GetByName(name string) (uint8, bool) {
	val, exist := platformsMap[name]
	return val, exist
}
