// Code generated by counterfeiter. DO NOT EDIT.
package internalfakes

import (
	"image"
	"io"
	"sync"

	"github.com/petewall/eink-radiator-image-source-image/v2/internal"
)

type FakeImageDecoder struct {
	Stub        func(io.Reader) (image.Image, error)
	mutex       sync.RWMutex
	argsForCall []struct {
		arg1 io.Reader
	}
	returns struct {
		result1 image.Image
		result2 error
	}
	returnsOnCall map[int]struct {
		result1 image.Image
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeImageDecoder) Spy(arg1 io.Reader) (image.Image, error) {
	fake.mutex.Lock()
	ret, specificReturn := fake.returnsOnCall[len(fake.argsForCall)]
	fake.argsForCall = append(fake.argsForCall, struct {
		arg1 io.Reader
	}{arg1})
	stub := fake.Stub
	returns := fake.returns
	fake.recordInvocation("ImageDecoder", []interface{}{arg1})
	fake.mutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return returns.result1, returns.result2
}

func (fake *FakeImageDecoder) CallCount() int {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return len(fake.argsForCall)
}

func (fake *FakeImageDecoder) Calls(stub func(io.Reader) (image.Image, error)) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = stub
}

func (fake *FakeImageDecoder) ArgsForCall(i int) io.Reader {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return fake.argsForCall[i].arg1
}

func (fake *FakeImageDecoder) Returns(result1 image.Image, result2 error) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = nil
	fake.returns = struct {
		result1 image.Image
		result2 error
	}{result1, result2}
}

func (fake *FakeImageDecoder) ReturnsOnCall(i int, result1 image.Image, result2 error) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = nil
	if fake.returnsOnCall == nil {
		fake.returnsOnCall = make(map[int]struct {
			result1 image.Image
			result2 error
		})
	}
	fake.returnsOnCall[i] = struct {
		result1 image.Image
		result2 error
	}{result1, result2}
}

func (fake *FakeImageDecoder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeImageDecoder) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ internal.ImageDecoder = new(FakeImageDecoder).Spy
