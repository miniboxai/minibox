package object_store

import "time"

type OptionFunc func(*ObjectStore)

type ObjectOption struct {
	Version string
	Tag     string
	Expire  time.Duration
}

type ObjectFunc func(opt *ObjectOption)

func Region(region string) OptionFunc {
	return func(objs *ObjectStore) {
		objs.Region = region
	}
}

func Endpoint(endpoint string) OptionFunc {
	return func(objs *ObjectStore) {
		objs.Endpoint = endpoint
	}
}

func StoreBucket(bucket string) OptionFunc {
	return func(objs *ObjectStore) {
		objs.StoreBucket = bucket
	}
}

func Root(pat string) OptionFunc {
	return func(objs *ObjectStore) {
		objs.Root = pat
	}
}

func Version(ver string) ObjectFunc {
	return func(opt *ObjectOption) {
		opt.Version = ver
	}
}
