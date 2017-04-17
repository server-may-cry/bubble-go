package application

import "github.com/server-may-cry/bubble-go/notification"

// VkEventChan channel for send vk events
var VkEventChan chan<- (notification.VkEvent)
