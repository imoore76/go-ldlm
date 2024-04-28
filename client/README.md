# go-ldlm Client

`go-ldlm/client` is a go library for communicating with an LDLM server (http://github.com/imoore76/go-ldlm).

## Installation

```go
// your go application
include "github.com/imoore76/go-ldlm/client"
```

The just run `go mod tidy`

## Usage

### Create a Client
A client takes a context and a pointer to a client.Config object.

```go
c := client.New(context.Background(), &client.Config{
    Address: "localhost:3144",
})
```

#### Config

client.Config members
| Name | Type | Description |
| :--- | :--- | :--- |
| `Address` | string  | host:port address of ldlm server |
| `Password` |  string  | Password to use for LDLM server |
| `NoAutoRefresh`  | bool  |  Don't automatically refresh locks in a background goroutine when a lock timeout is specified for a lock |
| `UseTls`    |  bool  |  Use TLS to connect to the server |
| `SkipVerify`  | bool  |  Don't verify the server's certificate |
| `CAFile`   | string |  File containing a CA certificate |
| `TlsCert`  |  string  | File containing a TLS certificate for this client |
| `TlsKey`   |  string  |  File containing a TLS key for this client |


### Basic Concepts

Locks in an LDLM server generally live until the client unlocks the lock or disconnects. If a client dies while holding a lock, the disconnection is detected and handled in LDLM by releasing the lock.

Depending on your LDLM server configuration, this feature may be disabled and `lockTimeoutSeconds` would be used to specify the maximum amount of time a lock can remain locked without being refreshed. The client will take care of refreshing locks in the background for you unless you've specified `NoAutoRefresh` in the client's options. Otherwise, you must periodically call `client.RefreshLock(...)` yourself &lt; the lock timeout interval.

To `Unlock()` or refresh a lock, you must use the lock key that was issued from the lock request's response. Using a Lock object's `Unlock()` method takes care of this for you. This is exemplified further in the following sections.

### Lock Object

`client.Lock` objects are returned from successful `Lock()` and `TryLock()` client methods. They have the following members:

| Name | Type | Description |
| :--- | :--- | :--- |
| `Name` | `string` | The name of the lock |
| `Key` | `string` | the key for the lock |
| `Locked` | `bool` | whether the was acquired or not |
| `Unlock()` | `func() (bool, error)` | method to unlock the lock |


### Lock

`Lock()` attempts to acquire a lock in LDLM. It will block until the lock is acquired or until `waitTimeoutSeconds` has elapsed (if specified).  If you have set `waitTimeoutSeconds` and the lock could not be acquired in that time, the returned error will be `client.ErrLockWaitTimeout`.

`Lock()` accepts the following arguments.

| Type | Description |
| :--- | :--- |
| `string` | Name of the lock to acquire |
| `*uint32` | Maximum amount of time to wait to acquire a lock. Use `nil` or `0` to wait indefinitely. |
| `*uint32` | The lock timeout. Use `nil` or `0` for no timeout.

It returns a `*Lock` and an `error`.



#### Examples

Simple lock
```go
lock, err = c.Lock("my-task", nil, nil)
if err != nill {
    // handle err
}
defer lock.Unlock()

// Perform task...

```

Wait timeout
```go
wait := uint32(30)
lock, err = c.Lock("my-task", &wait, nil)
if err != nill {
    // handle err
    if !errors.Is(err, client.ErrLockWaitTimeout) {
        // The error is not because the wait timeout was exceeded
    }
}
defer lock.Unlock()
```

### Try Lock
`TryLock()` attempts to acquire a lock and immediately returns; whether the lock was acquired or not. You must inspect the returned lock's `Locked` property to determine if it was acquired.

`TryLock()` accepts the following arguments.

| Type | Description |
| :--- | :--- |
| `string` | Name of the lock to acquire |
| `*uint32` | The lock timeout. Use `nil` or `0` for no timeout.

It returns a `*Lock` and an `error`.


#### Examples

Simple try lock
```go
lock, err = c.TryLock("my-task", nil)
if err != nill {
    // handle err
}
if !lock.Locked {
    // Something else is holding the lock
    return
}

defer lock.Unlock()
// Do work...

```


### Unlock
`Unlock()` unlocks the specified lock and stops any lock refresh job that may be associated with the lock. It must be passed the key that was issued when the lock was acquired. Using a different key will result in an error returned from LDLM and an error returned. Since an `Unlock()` method is available on Lock objects returned by the client, calling this directly should not be needed.

`Unlock()` accepts the following arguments.

| Type | Description |
| :--- | :--- |
| `string` | Name of the lock |
| `string` | Key for the lock |

It returns a `bool` indicating whether or not the lock was unlocked and an `error`.

#### Examples
Simple unlock
```go
unlocked, err := c.Unlock("my_task", lock.key)
```

### Refresh Lock
As explained in [Basic Concepts](#basic-concepts), you may specify a lock timeout using a `lockTimeoutSeconds` argument to any of the `*Lock*()` methods. When you do this, the client will refresh the lock in the background without you having to do anything. If, for some reason, you want to disable auto refresh (`NoAutoRefresh=true` in client Config), you will have to refresh the lock before it times out using the `RefreshLock()` method.

It takes the following arguments

| Type | Description |
| :--- | :--- |
| `string` | Name of the lock to acquire |
| `string` | The key for the lock |
| `uint32` | The new lock expiration timeout (or the same timeout if you'd like) |

It returns a `*Lock` and an `error`.

#### Examples
```go
lockTimeout := uint32(300)
lock, err = c.Lock("task1-lock", nil, &lockTimeout)

if err != nil {
    // handle err
}

defer lock.Unlock()

// do some work, then

if _, err := c.RefreshLock("task1-lock", lock.Key, 300); err != nil {
    panic(err)
}

// do some more work, then

if _, err := c.RefreshLock("task1-lock", lock.Key, 300); err != nil {
    panic(err)
}

// do some more work

```

## License

Apache 2.0; see [`LICENSE`](../LICENSE) for details.

## Contributing

See [`CONTRIBUTING.md`](../CONTRIBUTING.md) for details.

## Disclaimer

This project is not an official Google project. It is not supported by Google and Google specifically disclaims all warranties as to its quality, merchantability, or fitness for a particular purpose.
