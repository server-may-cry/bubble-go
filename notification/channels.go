package notification

// VkEvent struct for app2user event notification in Vk
type VkEvent struct {
	ExtID string
	Type  int
	Value int
}

// VkEventChan channel for sending app2user notification in VK platform
var VkEventChan chan (VkEvent)
