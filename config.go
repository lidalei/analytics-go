package analytics

import (
	"net/http"
	"time"

	"github.com/segmentio/backo-go"
	"github.com/xtgo/uuid"
)

// Instances of this type carry the different configuration options that may
// be set when instantiating a client.
//
// Each field's zero-value is either meaningful or interpreted as using the
// default value defined by the library.
type Config struct {

	// The endpoint to which the client connect and send their messages, set to
	// `DefaultEndpoint` by default.
	Endpoint string

	// The flushing interval of the client. Messages will be sent when they've
	// been queued up to the maximum batch size or when the flushing interval
	// timer triggers.
	Interval time.Duration

	// The HTTP transport used by the client, this allows an application to
	// redefine how requests are being sent at the HTTP level (for example,
	// to change the connection pooling policy).
	// If none is specified the client uses `http.DefaultTransport`.
	Transport http.Transport

	// The logger used by the client to output info or error messages when that
	// are generated by background operations.
	// If none is specified the client uses a standard logger that outputs to
	// `os.Stderr`.
	Logger Logger

	// The maximum number of messages that will be sent in one API call.
	// Messages will be sent when they've been queued up to the maximum batch
	// size or when the flushing interval timer triggers.
	// Note that the API will still enfore a 500KB limit on each HTTP request
	// which is independant from the number of embedded messages.
	BatchSize int

	// When set to true the client will send more frequent and detailed messages
	// to its logger.
	Verbose bool

	// A function called by the client to generate unique message identifiers.
	// The client uses a UUID generator if none is provided.
	UID func() string

	// A function called by the client to get the current time, `time.Now` is
	// used by default.
	Now func() time.Time

	// The retry policy used by the client to resend requests that have failed.
	// The function is called with how many times the operation has been retried
	// and is expected to return how long the client should wait before trying
	// again.
	// If not set the client will fallback to use a default retry policy.
	RetryAfter func(int) time.Duration
}

// This constant sets the default endpoint to which client instances send
// messages if none was explictly set.
const DefaultEndpoint = "https://api.segment.io"

// This constant sets the default flush interval used by client instances if
// none was explicitly set.
const DefaultInterval = 5 * time.Second

// This constant sets the default batch size used by client instances if none
// was explicitly set.
const DefaultBatchSize = 250

// Verifies that fields that don't have zero-values are set to valid values,
// returns an error describing the problem if a field was invalid.
func (c *Config) validate() error {
	if c.Interval < 0 {
		return ConfigError{
			Reason: "negative time intervals are not supported",
			Field:  "Interval",
			Value:  c.Interval,
		}
	}

	if c.BatchSize < 0 {
		return ConfigError{
			Reason: "negative batch sizes are not supported",
			Field:  "BatchSize",
			Value:  c.BatchSize,
		}
	}

	return nil
}

// Given a config object as argument the function will set all zero-values to
// their defaults and return the modified object.
func makeConfig(c Config) Config {
	if len(c.Endpoint) == 0 {
		c.Endpoint = DefaultEndpoint
	}

	if c.Interval == 0 {
		c.Interval = DefaultInterval
	}

	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}

	if c.Logger == nil {
		c.Logger = newDefaultLogger()
	}

	if c.BatchSize == 0 {
		c.BatchSize = DefaultBatchSize
	}

	if c.UID == nil {
		c.UID = uid
	}

	if c.Now == nil {
		c.Now = time.Now
	}

	if c.RetryAfter == nil {
		c.RetryAfter = backo.DefaultBacko().Duration
	}

	return c
}

// This function returns a string representation of a UUID, it's the default
// function used for generating unique IDs.
func uid() string {
	return uuid.NewRandom().String()
}
