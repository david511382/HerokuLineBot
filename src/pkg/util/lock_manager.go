package util

import (
	"sync"
)

type lockRequest struct {
	lockerNames []string
	freeFnCh    chan func()
}

type LockManager struct {
	lockerRegisterMap map[string]*sync.RWMutex
	lockCh            chan *lockRequest
}

func NewLockManager(lockerNames ...string) *LockManager {
	r := &LockManager{
		lockerRegisterMap: make(map[string]*sync.RWMutex),
		lockCh:            make(chan *lockRequest, 1),
	}
	for _, name := range lockerNames {
		r.lockerRegisterMap[name] = &sync.RWMutex{}
	}
	go r.handleLockReqs()
	return r
}

func (l *LockManager) Lock(lockerNames ...string) (freeFnCh chan func()) {
	freeFnCh = make(chan func(), 1)
	go l.queueLock(freeFnCh, lockerNames...)
	return
}

func (l *LockManager) queueLock(freeFnCh chan func(), lockerNames ...string) {
	l.lockCh <- &lockRequest{
		lockerNames: lockerNames,
		freeFnCh:    freeFnCh,
	}
}

func (l *LockManager) handleLockReqs() {
	for {
		reqs := <-l.lockCh
		go func(reqs *lockRequest) {
			var unLockFns []func()
			for _, v := range reqs.lockerNames {
				locker := l.lockerRegisterMap[v]
				if locker == nil {
					continue
				}

				success := locker.TryLock()
				if !success {
					// free all
					for _, v := range unLockFns {
						v()
					}
					// requeue
					l.queueLock(reqs.freeFnCh, reqs.lockerNames...)
					return
				}
				unLockFns = append(unLockFns, locker.Unlock)
			}

			reqs.freeFnCh <- func() {
				for _, v := range unLockFns {
					v()
				}
			}
		}(reqs)
	}
}
