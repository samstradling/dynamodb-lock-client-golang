package lockclient

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/atomic"
)

// DynamoDBLockClient describes the fields for a lock client
type DynamoDBLockClient struct {
	LockName        string
	LeaseDuration   time.Duration
	HeartbeatPeriod time.Duration
	TableName       string
	Identifier      string
	Client          dynamodbiface.DynamoDBAPI
	lockID          string
	sendHeartbeats  atomic.Bool
	lockError       error
}
