/**
 * @Author: ChenJunJi
 * @Desc:
 * @Date: 2021/8/27 17:11
 */

package lock

import "sync"

type LockSt struct {
	sync.RWMutex
}

func NewWithUsageMu(internalMu sync.Locker, onUsageAdd func(), onUsageSub func()) *WithUsageMu {
	return &WithUsageMu{
		onUsageAdd: onUsageAdd,
		onUsageSub: onUsageSub,
		internalMu: internalMu,
	}
}

type WithUsageMu struct {
	outerMu    sync.RWMutex
	internalMu sync.Locker
	UsageNum   uint32
	onUsageSub func()
	onUsageAdd func()
}

func (m *WithUsageMu) Lock() {
	m.internalMu.Lock()
	m.UsageNum++
	m.onUsageAdd()
	m.internalMu.Unlock()

	m.outerMu.Lock()
}

func (m *WithUsageMu) Unlock() {
	m.internalMu.Lock()
	m.UsageNum--
	m.onUsageSub()
	m.internalMu.Unlock()

	m.outerMu.Unlock()
}

func (m *WithUsageMu) RLock() {
	m.internalMu.Lock()
	m.UsageNum++
	m.onUsageAdd()
	m.internalMu.Unlock()

	m.outerMu.RLock()
}

func (m *WithUsageMu) RUnlock() {
	m.internalMu.Lock()
	m.UsageNum--
	m.onUsageSub()
	m.internalMu.Unlock()

	m.outerMu.RUnlock()
}

type MulElemMuFactory struct {
	elemMuMap sync.Map
	opMapMu   sync.Mutex
}

func NewMulElemMuFactory() *MulElemMuFactory {
	return &MulElemMuFactory{}
}

func (m *MulElemMuFactory) MakeOrGetSpecElemMu(elem interface{}) *WithUsageMu {
	mu, ok := m.elemMuMap.Load(elem)
	if !ok {
		mu = NewWithUsageMu(
			&m.opMapMu,
			func() {
				// save this lock while any thread used the handler lock
				if _, ok := m.elemMuMap.Load(elem); ok {
					return
				}
				m.elemMuMap.Store(elem, mu)
			},
			func() {
				// remove this lock from map if this lock might have no owner
				if mu == nil {
					m.elemMuMap.Delete(elem)
					return
				}
				lmu, ok := mu.(*WithUsageMu)
				if !ok {
					m.elemMuMap.Delete(elem)
					return
				}
				// remove this lock from map if this lock might have no owner
				if lmu.UsageNum > 0 {
					return
				}
				m.elemMuMap.Delete(elem)
			},
		)

		m.elemMuMap.Store(elem, mu)
	}
	return mu.(*WithUsageMu)
}
