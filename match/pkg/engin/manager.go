package engin

import (
	"sync"
)

// 多市場撮合引擎管理器
type EngineManager struct {
	mu      sync.RWMutex
	engines map[string]*Engine
}

// 初始化 EngineManager，內建三個市場
func NewEngineManager(engines map[string]*Engine) *EngineManager {
	mgr := &EngineManager{
		engines: engines,
	}
	return mgr
}

// 根據 symbol 取得或創建撮合引擎（具備 thread-safe）
func (m *EngineManager) GetOrCreateEngine(symbol string) *Engine {
	m.mu.RLock()
	engine, ok := m.engines[symbol]
	m.mu.RUnlock()
	if ok {
		return engine
	}

	// 若不存在則創建
	m.mu.Lock()
	defer m.mu.Unlock()
	// double-check
	if engine, ok := m.engines[symbol]; ok {
		return engine
	}
	engine = NewEngine()
	m.engines[symbol] = engine
	return engine
}
