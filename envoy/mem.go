package envoy

import (
	"sync"
	"unsafe"
)

var memManager memoryManager

type (
	// memoryManager manages the heap allocated context objects.
	// It is used to pin the context objects to the heap to avoid them being garbage collected by the Go runtime.
	//
	// TODO: shard the linked list of HttpContexts to reduce contention.
	//
	// TODO: is this really necessary? Pinning a pointer to the interface might work? e.g.
	// 	...
	//   pinner := runtime.Pinner{}
	//   wrapper := &pinedHttpContext{ctx: ctx}
	//   pinn.Pinned(wrapper)
	// 	...
	//  does this work even when the data inside the interface contains pointers?
	memoryManager struct {
		// moduleContext holds the module context.
		moduleContext ModuleContext
		// httpContexts holds a linked lists of HttpContext.
		httpContexts      *pinedHttpContext
		httpContextsMutex sync.Mutex
	}
	// pinnedHttpContext holds a pinned HttpContext managed by the memory manager.
	pinedHttpContext struct {
		ctx        HttpContext
		next, prev *pinedHttpContext
	}
)

// pinModuleContext pins the module context to the memory manager.
func (m *memoryManager) pinModuleContext(ctx ModuleContext) {
	m.moduleContext = ctx
}

func (m *memoryManager) unwrapPinnedModuleContext() *ModuleContext {
	return &m.moduleContext
}

func (m *memoryManager) pinHttpContext(ctx HttpContext) *pinedHttpContext {
	m.httpContextsMutex.Lock()
	defer m.httpContextsMutex.Unlock()
	item := &pinedHttpContext{ctx: ctx, next: m.httpContexts, prev: nil}
	if m.httpContexts != nil {
		m.httpContexts.prev = item
	}
	m.httpContexts = item
	return item
}

func (m *memoryManager) removeHttpContext(ctx *pinedHttpContext) {
	m.httpContextsMutex.Lock()
	defer m.httpContextsMutex.Unlock()
	if ctx.prev != nil {
		ctx.prev.next = ctx.next
	} else {
		m.httpContexts = ctx.next
	}
	if ctx.next != nil {
		ctx.next.prev = ctx.prev
	}
}

func unwrapRawPinHttpContext(rawPingedHttpContext uintptr) HttpContext {
	return (*pinedHttpContext)(unsafe.Pointer(rawPingedHttpContext)).ctx
}
