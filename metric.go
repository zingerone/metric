package metric

import (
	"fmt"
	"time"

	vstatsd "github.com/DataDog/datadog-go/statsd"
)

var statsd *vstatsd.Client

// Option this for options
type Option func(*Options) error

// Options this for setting options
type Options struct {
	serviceName         string
	env                 string
	maxMessageOnPayload int
}

func resolveOptions(opt []Option) (*Options, error) {
	o := &Options{
		serviceName:         "unknown",
		env:                 "dev",
		maxMessageOnPayload: 100,
	}
	for _, option := range opt {
		err := option(o)
		if err != nil {
			return o, err
		}
	}
	return o, nil

}
func ServiceName(serviceName string) Option {
	return func(o *Options) error {
		o.serviceName = serviceName
		return nil
	}
}

func SetEnv(env string) Option {
	return func(o *Options) error {
		o.env = env
		return nil
	}
}

func SetMaxMessageOnPayload(maxProc int) Option {
	return func(o *Options) error {
		o.maxMessageOnPayload = maxProc
		return nil
	}
}

// this function for set global tracer
func SetGlobal(host, port string, opts ...Option) error {
	opt, _ := resolveOptions(opts)
	addr := fmt.Sprintf("%s:%s", host, port)
	tempStatsd, err := vstatsd.New(addr,
		vstatsd.WithTags([]string{"env:" + opt.env, "service:" + opt.serviceName}),
		vstatsd.WithMaxMessagesPerPayload(opt.maxMessageOnPayload),
	)

	if err != nil {
		return err
	}
	statsd = tempStatsd
	return nil
}

// Close this for ensure all metric deliver
func Close() error {
	return statsd.Close()
}

// Tracer struct metric
type Tracer struct {
	name  string
	start time.Time
	tags  []string
}

// AddTags this function for add tags
func (t *Tracer) AddTags(tags Tags) *Tracer {
	temp := tagsToArr(tags)
	if len(temp) > 0 {
		t.tags = append(t.tags, temp...)
	}
	return t
}

// Stop this for stop tracer
func (t *Tracer) Stop(err error) error {
	currentTime := time.Now()
	diff := currentTime.Sub(t.start)
	if err != nil {
		t.tags = append(t.tags, fmt.Sprintf("err:%s", err.Error()))

	}
	return statsd.Histogram(t.name, diff.Seconds(), t.tags, 1.0)
}

// Tags this for add tags
type Tags map[string]interface{}

// NewTrace this for new tracer
func NewTrace(name string) (Tracer, error) {
	var (
		sTags []string
	)

	return Tracer{
		name:  name,
		start: time.Now(),
		tags:  sTags,
	}, nil
}

// tagsToArr convert tags to Arr
func tagsToArr(tags Tags) []string {
	var sTags []string
	for k, v := range tags {
		sTags = append(sTags, fmt.Sprintf("%s:%v", k, v))
	}
	return sTags
}
