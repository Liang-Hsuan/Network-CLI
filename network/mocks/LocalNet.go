// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import network "../"

// LocalNet is an autogenerated mock type for the LocalNet type
type LocalNet struct {
	mock.Mock
}

// Dial provides a mock function with given fields: protocol, address
func (_m *LocalNet) Dial(protocol string, address string) (network.LocalNetConn, error) {
	ret := _m.Called(protocol, address)

	var r0 network.LocalNetConn
	if rf, ok := ret.Get(0).(func(string, string) network.LocalNetConn); ok {
		r0 = rf(protocol, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.LocalNetConn)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(protocol, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Listen provides a mock function with given fields: protocol, address
func (_m *LocalNet) Listen(protocol string, address string) (network.LocalNetListener, error) {
	ret := _m.Called(protocol, address)

	var r0 network.LocalNetListener
	if rf, ok := ret.Get(0).(func(string, string) network.LocalNetListener); ok {
		r0 = rf(protocol, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.LocalNetListener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(protocol, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
