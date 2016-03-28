package cmap

import (
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

type innerChan struct {
	q       []aValue
	sendIdx uint32
	recvIdx uint32
	slen    uint32
	rlen    uint32
	die     uint32
}

type Chan struct {
	*innerChan
}

func newKeyValueChan() Chan {
	return newSizeKeyValueChan(1)
}

func newSizeKeyValueChan(sz int) Chan {
	if sz < 1 {
		panic("sz < 1")
	}
	return Chan{&innerChan{
		q:       make([]aValue, sz),
		sendIdx: ^uint32(0),
		recvIdx: ^uint32(0),
	}}
}

func (ch Chan) Send(v *KeyValue, block bool) bool {
	if !block && ch.Len() == ch.Cap() {
		return false
	}
	ln, cnt := uint32(len(ch.q)), uint32(0)
	for !ch.Closed() {
		if ch.Len() == ch.Cap() {
			if !block {
				return false
			}
			runtime.Gosched()
			continue
		}
		i := atomic.AddUint32(&ch.sendIdx, 1)
		if ch.q[i%ln].CompareAndSwapIfNil(v) {
			atomic.AddUint32(&ch.slen, 1)
			return true
		}
		if block {
			if i%250 == 0 {
				pause(1)
			}
		} else if cnt++; cnt == ln {
			break
		}
		runtime.Gosched()
	}
	return false
}

func (ch Chan) Recv(block bool) (*KeyValue, bool) {
	if !block && ch.Len() == 0 { // fast path
		return zeroValue, false
	}
	ln, cnt := uint32(len(ch.q)), uint32(0)
	for !ch.Closed() || ch.Len() > 0 {
		if ch.Len() == 0 {
			if !block {
				return zeroValue, false
			}
			runtime.Gosched()
			continue
		}
		i := atomic.AddUint32(&ch.recvIdx, 1)
		if v, ok := ch.q[i%ln].SwapWithNil(); ok {
			atomic.AddUint32(&ch.rlen, 1)
			return v, true
		}
		if block {
			if i%250 == 0 {
				pause(1)
			}
		} else if cnt++; cnt == ln {
			break
		}
		runtime.Gosched()
	}
	return zeroValue, false
}

func (ch Chan) SendOnly() SendOnly { return SendOnly{ch} }

func (ch Chan) RecvOnly() RecvOnly { return RecvOnly{ch} }

func (ch Chan) Close() { atomic.StoreUint32(&ch.die, 1) }

func (ch Chan) Closed() bool { return atomic.LoadUint32(&ch.die) == 1 }

func (ch Chan) Cap() int { return len(ch.q) }

func (ch Chan) Len() int { return int(atomic.LoadUint32(&ch.slen) - atomic.LoadUint32(&ch.rlen)) }

func SelectSend(block bool, v *KeyValue, chans ...Sender) bool {
	for {
		for i := range chans {
			if ok := chans[i].Send(v, false); ok {
				return ok
			}
		}
		if !block {
			return false
		}
		pause(1)
	}
}

func SelectRecv(block bool, chans ...Receiver) (*KeyValue, bool) {
	for {
		for i := range chans {
			if v, ok := chans[i].Recv(false); ok {
				return v, ok
			}
		}
		if !block {
			return zeroValue, false
		}
		pause(1)
	}
}

type SendOnly struct{ c Chan }

func (so SendOnly) Send(v *KeyValue, block bool) bool { return so.c.Send(v, block) }

type Sender interface {
	Send(v *KeyValue, block bool) bool
}

type RecvOnly struct{ c Chan }

func (ro RecvOnly) Recv(block bool) (*KeyValue, bool) { return ro.c.Recv(block) }

type Receiver interface {
	Recv(block bool) (*KeyValue, bool)
}

func pause(p time.Duration) { time.Sleep(time.Millisecond * p) }

var (
	_ Sender   = (*Chan)(nil)
	_ Sender   = (*SendOnly)(nil)
	_ Receiver = (*Chan)(nil)
	_ Receiver = (*RecvOnly)(nil)
)

var zeroValue *KeyValue

type aValue struct {
	v *KeyValue
}

func (a *aValue) CompareAndSwapIfNil(newVal *KeyValue) bool {
	x := unsafe.Pointer(&a.v)
	return atomic.CompareAndSwapPointer((*unsafe.Pointer)(atomic.LoadPointer(&x)), nil, unsafe.Pointer(&newVal))
}

func (a *aValue) SwapWithNil() (*KeyValue, bool) {
	x := unsafe.Pointer(&a.v)
	if v := atomic.SwapPointer((*unsafe.Pointer)(atomic.LoadPointer(&x)), nil); v != nil {
		return *(**KeyValue)(v), true
	}
	return zeroValue, false
}
