package lockclient

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func int64New(x int64) *int64 {
	return &x
}

type MockedDynamoClient struct {
	dynamodb.DynamoDB
	ScanResp       dynamodb.ScanOutput
	Err            error
	FailAfterFirst bool
	PutItemCount   int
}

func (d *MockedDynamoClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return &d.ScanResp, d.Err
}

func (d *MockedDynamoClient) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	d.PutItemCount++
	if d.FailAfterFirst && d.PutItemCount > 1 {
		return &dynamodb.PutItemOutput{}, errors.New("An error occured")
	}
	return &dynamodb.PutItemOutput{}, d.Err
}

func (d *MockedDynamoClient) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, d.Err
}

func TestHasLock(t *testing.T) {
	cases := []struct {
		Resp    dynamodb.ScanOutput
		err     error
		GotLock bool
		name    string
	}{
		{
			name: "Test successfull write returns true",
			Resp: dynamodb.ScanOutput{
				Count: int64New(1),
			},
			err:     nil,
			GotLock: true,
		},
		{
			name: "Test unsuccessfull write returns false",
			Resp: dynamodb.ScanOutput{
				Count: int64New(0),
			},
			err:     nil,
			GotLock: false,
		},
		{
			name: "Test client returns an error returns false",
			Resp: dynamodb.ScanOutput{
				Count: int64New(1),
			},
			err:     errors.New("An error occured"),
			GotLock: false,
		},
	}

	for i, c := range cases {
		d := &DynamoDBLockClient{
			LockName:        "my-unique-lock-name",
			LeaseDuration:   5000 * time.Millisecond,
			HeartbeatPeriod: 1000 * time.Millisecond,
			TableName:       "LockTable",
			Client:          &MockedDynamoClient{ScanResp: c.Resp, Err: c.err},
		}
		hasLock, err := d.HasLock()
		if hasLock != c.GotLock {
			t.Fatalf("HasLock() Test %d (%s) failed: unexpected response: %t, expecting: %t", i, c.name, hasLock, c.GotLock)
		}
		if err != c.err {
			t.Fatalf("HasLock() Test %d (%s) failed: unexpected error passed on: %s expecting: %s", i, c.name, err, c.err)
		}
	}
}

func TestGetLock(t *testing.T) {
	cases := []struct {
		err    error
		name   string
		result bool
	}{
		{
			name:   "Test client passes on error: nil",
			err:    nil,
			result: true,
		},
		{
			name:   "Test client passes on error: an error occured",
			err:    errors.New("An error occured"),
			result: false,
		},
	}

	for i, c := range cases {
		d := &DynamoDBLockClient{
			LockName:        "my-unique-lock-name",
			LeaseDuration:   5000 * time.Millisecond,
			HeartbeatPeriod: 1000 * time.Millisecond,
			TableName:       "LockTable",
			Client:          &MockedDynamoClient{Err: c.err},
		}
		result, err := d.GetLock()
		if err != c.err {
			t.Fatalf("HasLock() Test %d (%s) failed: unexpected error passed on: %s expecting: %s", i, c.name, err, c.err)
		}
		if result != c.result {
			t.Fatalf("HasLock() Test %d (%s) failed: unexpected result: %t expecting: %t", i, c.name, result, c.result)
		}
	}
}

func TestRemoveLock(t *testing.T) {
	cases := []struct {
		err  error
		name string
	}{
		{
			name: "Test client passes on error: nil",
			err:  nil,
		},
		{
			name: "Test client passes on error: an error occured",
			err:  errors.New("An error occured"),
		},
	}

	for i, c := range cases {
		d := &DynamoDBLockClient{
			LockName:        "my-unique-lock-name",
			LeaseDuration:   5000 * time.Millisecond,
			HeartbeatPeriod: 1000 * time.Millisecond,
			TableName:       "LockTable",
			Client:          &MockedDynamoClient{Err: c.err},
		}
		err := d.RemoveLock()
		if err != c.err {
			t.Fatalf("HasLock() Test %d (%s) failed: unexpected error passed on: %s expecting: %s", i, c.name, err, c.err)
		}
	}
}

func TestStopHeartbeat(t *testing.T) {

	MockedClient := MockedDynamoClient{Err: nil}
	d := &DynamoDBLockClient{
		LockName:        "my-unique-lock-name",
		LeaseDuration:   50 * time.Millisecond,
		HeartbeatPeriod: 10 * time.Millisecond,
		TableName:       "LockTable",
		Client:          &MockedClient,
	}

	d.GetLock()

	time.Sleep(50 * time.Millisecond)

	if MockedClient.PutItemCount == 0 {
		t.Fatalf("Client did not start sending heartbeats")
	}

	d.StopHeartbeat()
	MockedClient.PutItemCount = 0

	time.Sleep(50 * time.Millisecond)

	if MockedClient.PutItemCount > 0 {
		t.Fatalf("Client did not stop sending heartbeats")
	}
}

func TestHeartbeatFails(t *testing.T) {

	d := &DynamoDBLockClient{
		LockName:        "my-unique-lock-name",
		LeaseDuration:   50 * time.Millisecond,
		HeartbeatPeriod: 10 * time.Millisecond,
		TableName:       "LockTable",
		Client:          &MockedDynamoClient{Err: nil, FailAfterFirst: true},
	}

	d.GetLock()

	time.Sleep(50 * time.Millisecond)

	if d.LockError() == nil {
		t.Fatalf("Client failed to notice failed heartbeat")
	}
}

func IsALockClient(lc LockClient) func() (bool, error) {
	return lc.HasLock
}

func TestDynamoDBLockClientSatisfiesRequirementsToBeALockClient(t *testing.T) {
	IsALockClient(&DynamoDBLockClient{}) // Gives a compile time error if it doesn't satisfy the interface
}
