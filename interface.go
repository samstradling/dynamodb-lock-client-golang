package lockclient

// LockClient interface names the method signatures for implementing a DynamoDB Lock Client
type LockClient interface {
	GetLock() (bool, error)
	HasLock() (bool, error)
	StopHeartbeat()
	RemoveLock() error
}
