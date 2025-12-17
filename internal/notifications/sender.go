package notifications

func Send(n Notification) {
	select {
	case NotificationsChannel <- n:
	default:
	}
}
