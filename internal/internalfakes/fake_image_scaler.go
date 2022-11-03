// Code generated by counterfeiter. DO NOT EDIT.
package internalfakes

import (
	"image"
	"image/draw"
	"sync"

	"github.com/petewall/eink-radiator-image-source-image/internal"
	drawa "golang.org/x/image/draw"
)

type FakeImageScaler struct {
	Stub        func(draw.Image, image.Rectangle, image.Image, image.Rectangle, draw.Op, *drawa.Options)
	mutex       sync.RWMutex
	argsForCall []struct {
		arg1 draw.Image
		arg2 image.Rectangle
		arg3 image.Image
		arg4 image.Rectangle
		arg5 draw.Op
		arg6 *drawa.Options
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeImageScaler) Spy(arg1 draw.Image, arg2 image.Rectangle, arg3 image.Image, arg4 image.Rectangle, arg5 draw.Op, arg6 *drawa.Options) {
	fake.mutex.Lock()
	fake.argsForCall = append(fake.argsForCall, struct {
		arg1 draw.Image
		arg2 image.Rectangle
		arg3 image.Image
		arg4 image.Rectangle
		arg5 draw.Op
		arg6 *drawa.Options
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	stub := fake.Stub
	fake.recordInvocation("ImageScaler", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.mutex.Unlock()
	if stub != nil {
		fake.Stub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
}

func (fake *FakeImageScaler) CallCount() int {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return len(fake.argsForCall)
}

func (fake *FakeImageScaler) Calls(stub func(draw.Image, image.Rectangle, image.Image, image.Rectangle, draw.Op, *drawa.Options)) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = stub
}

func (fake *FakeImageScaler) ArgsForCall(i int) (draw.Image, image.Rectangle, image.Image, image.Rectangle, draw.Op, *drawa.Options) {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return fake.argsForCall[i].arg1, fake.argsForCall[i].arg2, fake.argsForCall[i].arg3, fake.argsForCall[i].arg4, fake.argsForCall[i].arg5, fake.argsForCall[i].arg6
}

func (fake *FakeImageScaler) Invocations() map[string][][]interface{} {
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

func (fake *FakeImageScaler) recordInvocation(key string, args []interface{}) {
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

var _ internal.ImageScaler = new(FakeImageScaler).Spy
