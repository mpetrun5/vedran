// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Backoff is an autogenerated mock type for the Backoff type
type Backoff struct {
	mock.Mock
}

// NextBackOff provides a mock function with given fields:
func (_m *Backoff) NextBackOff() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// Reset provides a mock function with given fields:
func (_m *Backoff) Reset() {
	_m.Called()
}
