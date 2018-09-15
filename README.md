# DynamoDB Lock Client - Golang

[![Build Status](https://travis-ci.org/samstradling/dynamodb-lock-client-golang.svg?branch=master)](https://travis-ci.org/samstradling/dynamodb-lock-client-golang)
[![Coverage Status](https://coveralls.io/repos/github/samstradling/dynamodb-lock-client-golang/badge.svg?branch=master)](https://coveralls.io/github/samstradling/dynamodb-lock-client-golang?branch=master)

A golang 1.x implementation of [awslabs/dynamodb-lock-client](https://github.com/awslabs/dynamodb-lock-client).

>The Amazon DynamoDB Lock Client is a general purpose distributed locking library built for DynamoDB. The DynamoDB Lock Client supports both fine-grained and coarse-grained locking as the lock keys can be any arbitrary string, up to a certain length.<sup>[[1]](https://github.com/awslabs/dynamodb-lock-client)</sup>

## Table of Contents

* [Usage](#usage)
* [Documentation](#documentation)
    * [lockclient.DynamoDBLockClient](#lockclientdynamodblockclient)
    * [GetLock()](#getlock)
    * [StopHeartbeat()](#stopheartbeat)
    * [RemoveLock()](#removelock)
    * [HasLock()](#haslock)
    * [SetLevel()](#setlevel)
* [Contributing](#contributing)

## Usage

To pull dependencies:

```bash
go get github.com/samstradling/dynamodb-lock-client-golang
```

For a full example see [examples/basic-case/main.go](./examples/basic-case/main.go).

```go
// Setup AWS region, and, for local development only, the endpoint
config := &aws.Config{
        Region:   aws.String("eu-west-1"),
        Endpoint: aws.String("http://localhost:8000"),
    }
sess := session.Must(session.NewSession(config))
// Setup the lock client config
lockClient := &lockclient.DynamoDBLockClient{
    LockName:        "my-unique-lock-name",
    LeaseDuration:   5000 * time.Millisecond,
    HeartbeatPeriod: 1000 * time.Millisecond,
    TableName:       "LockTable",
    Client:          dynamodb.New(sess),
}

for true {
    err := lockClient.GetLock() // Try to get the lock
    if err != nil {
        break
    }
    // Keep trying every second
    time.Sleep(1 * time.Second)
}
defer lockClient.RemoveLock() // remove the lock once done

for i := 0; i < 10; i++ {
    _, err := lockClient.HasLock() // check if lock is still valid
    if err != nil {
        logrus.Debug(err)
    }
    time.Sleep(1000 * time.Millisecond) // super important action that requires a global lock
}
```

## Documentation

```golang
import "github.com/samstradling/dynamodb-lock-client-golang"
```

To use, create a dynamo table with a hash key named `key`.

### lockclient.DynamoDBLockClient

DynamoDBLockClient creates a new lock client.

```golang
lockClient := &lockclient.DynamoDBLockClient{}
```

[`DynamoDBLockClient` fields](./data.go):

| Field             | Type          | Description |
|-------------------|---------------|-|
| `LockName`        | `string`      | A unique string that represents the lock required. A client will only compete with other clients presenting the same `LockName` string. |
| `LeaseDuration`   | `Duration`    | The time to lease the lock, for the initial request and every renewal. |
| `HeartbeatPeriod` | `Duration`    | The time between renewals of the lock. At maximum should be `LeaseDuration / 2`. |
| `TableName`       | `string`      | DynamoDB table name. |
| `Identifier`      | `string`      | Used for uniquely identifying the client. If unset a random one is generated using `uuid.NewRandom()` |
| `Client`          | `*DynamoDB`   | DynamoDB client used for the connection. |

### GetLock()

Request a new lock, returns true on success.

```golang
func (d *DynamoDBLockClient) GetLock() (bool, error)
```

### StopHeartbeat()

Stop sending lock heartbeats, any existing locks will not be removed until they expire.

```golang
func (d *DynamoDBLockClient) StopHeartbeat()
```

### RemoveLock()

Request the removal of any existing locks.

```golang
func (d *DynamoDBLockClient) RemoveLock() error
```

### HasLock()

Returns `true` if the existing lock is valid.

```golang
func (d *DynamoDBLockClient) HasLock() (bool, error)
```

### SetLevel()

This library uses [github.com/sirupsen/logrus](github.com/sirupsen/logrus). To set the log level use:

```golang
logrus.SetLevel(logrus.DebugLevel)
logrus.SetLevel(logrus.InfoLevel)
logrus.SetLevel(logrus.ErrorLevel)
```

## Contributing

To develop locally you'll probably want to run a local version of dynamodb. This can be done with the following docker command:

```bash
docker run -p 8000:8000 -d amazon/dynamodb-local
```

The examples directory contains a script to create a local dynamodb table:

```bash
go run examples/create-local-table/main.go
```

Contibutions welcomed.
