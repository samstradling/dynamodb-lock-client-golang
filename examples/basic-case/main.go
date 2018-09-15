package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sirupsen/logrus"

	"github.com/samstradling/dynamodb-lock-client-golang"
)

func actionThatRequiresLock() {

	logrus.SetLevel(logrus.DebugLevel)

	config := &aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String("http://localhost:8000"),
	}
	sess := session.Must(session.NewSession(config))
	lockClient := &lockclient.DynamoDBLockClient{
		LockName:        "my-unique-lock-name",
		LeaseDuration:   5000 * time.Millisecond,
		HeartbeatPeriod: 1000 * time.Millisecond,
		TableName:       "LockTable",
		Client:          dynamodb.New(sess),
	}

	for true {
		result, err := lockClient.GetLock()
		if result {
			break
		}
		logrus.Debug("Unabe to get lock: ", err)
		time.Sleep(1 * time.Second)
	}
	defer lockClient.RemoveLock()

	logrus.Debug("Successfully got lock!")

	// Do action that requires lock
	for i := 0; i < 10; i++ {
		_, err := lockClient.HasLock() // check if lock is still valid
		if err != nil {
			logrus.Debug(err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

}

func main() {

	actionThatRequiresLock()
	time.Sleep(5 * time.Second) // This is just here to show that heartbeats have stopped

}
