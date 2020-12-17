package watcher

// Notification encapsulates the data we send as the notification.
type Notification struct {
	UUID    string
	OldData string
	NewData string
}

// Watcher interface allows extending to other databases.
type Watcher interface {
	// Create the audit table within the database
	CreateAuditTable() error
	// Create a trigger to write changes into the audit table.
	CreateTrigger() error
	// Watch starts up the watcher that watches for changes then sends them to the returned channel.
	Watch() chan Notification
}
