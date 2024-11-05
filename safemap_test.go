package safemap

import (
	"sync"
	"testing"
)

func TestSafeMap_StoreAndLoad(t *testing.T) {
	sm := New[int, string]()
	sm.Store(1, "one")

	value, ok := sm.Load(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}
}

func TestSafeMap_Exists(t *testing.T) {
	sm := New[int, string]()
	sm.Store(1, "one")

	if !sm.Exists(1) {
		t.Errorf("expected key 1 to exist")
	}

	if sm.Exists(2) {
		t.Errorf("expected key 2 to not exist")
	}
}

func TestSafeMap_Delete(t *testing.T) {
	sm := New[int, string]()
	sm.Store(1, "one")
	sm.Delete(1)

	if sm.Exists(1) {
		t.Errorf("expected key 1 to be deleted")
	}
}

func TestSafeMap_Add(t *testing.T) {
	sm := New[int, int]()
	sm.Store(1, 10)
	sm.Add(1, 5)

	value, ok := sm.Load(1)
	if !ok || value != 15 {
		t.Errorf("expected 15, got '%v'", value)
	}
}

func TestSafeMap_Sub(t *testing.T) {
	sm := New[int, int]()
	sm.Store(1, 10)
	sm.Sub(1, 5)

	value, ok := sm.Load(1)
	if !ok || value != 5 {
		t.Errorf("expected 5, got '%v'", value)
	}
}

func TestSafeMap_LoadAndDelete(t *testing.T) {
	sm := New[int, string]()
	sm.Store(1, "one")

	value, ok := sm.LoadAndDelete(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}

	if sm.Exists(1) {
		t.Errorf("expected key 1 to be deleted")
	}
}

func TestSafeMap_LoadOrStore(t *testing.T) {
	sm := New[int, string]()
	value, loaded := sm.LoadOrStore(1, "one")
	if loaded || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}

	value, loaded = sm.LoadOrStore(1, "two")
	if !loaded || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}
}

func TestSafeMap_Swap(t *testing.T) {
	sm := New[int, string]()
	sm.Store(1, "one")

	previous, loaded := sm.Swap(1, "two")
	if !loaded || previous != "one" {
		t.Errorf("expected 'one', got '%v'", previous)
	}

	value, ok := sm.Load(1)
	if !ok || value != "two" {
		t.Errorf("expected 'two', got '%v'", value)
	}
}

func TestSafeMap_Range(t *testing.T) {
	sm := New[int, string]()
	sm.Store(1, "one")
	sm.Store(2, "two")

	keys := make(map[int]bool)
	sm.Range(func(k int, v string) bool {
		keys[k] = true
		return true
	})

	if len(keys) != 2 || !keys[1] || !keys[2] {
		t.Errorf("expected keys 1 and 2, got '%v'", keys)
	}
}

func TestSafeMap_Clear(t *testing.T) {
	sm := New[int, string]()
	sm.Store(1, "one")
	sm.Clear()

	if sm.Len() != 0 {
		t.Errorf("expected length 0, got '%v'", sm.Len())
	}
}

func TestSafeMap_KeysInRange(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")
	sm.Store(2, "two")
	sm.Store(3, "three")

	keys := sm.KeysInRange(1, 3)
	if len(keys) != 2 || keys[0] != 1 || keys[1] != 2 {
		t.Errorf("expected keys 1 and 2, got '%v'", keys)
	}
}
func TestSafeMap_StoreAndLoad_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")

	value, ok := sm.Load(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}
}

func TestSafeMap_Exists_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")

	if !sm.Exists(1) {
		t.Errorf("expected key 1 to exist")
	}

	if sm.Exists(2) {
		t.Errorf("expected key 2 to not exist")
	}
}

func TestSafeMap_Delete_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")
	sm.Delete(1)

	if sm.Exists(1) {
		t.Errorf("expected key 1 to be deleted")
	}
}

func TestSafeMap_Add_Ordered(t *testing.T) {
	sm := New[int, int](Ordered)
	sm.Store(1, 10)
	sm.Add(1, 5)

	value, ok := sm.Load(1)
	if !ok || value != 15 {
		t.Errorf("expected 15, got '%v'", value)
	}
}

func TestSafeMap_Sub_Ordered(t *testing.T) {
	sm := New[int, int](Ordered)
	sm.Store(1, 10)
	sm.Sub(1, 5)

	value, ok := sm.Load(1)
	if !ok || value != 5 {
		t.Errorf("expected 5, got '%v'", value)
	}
}

func TestSafeMap_LoadAndDelete_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")

	value, ok := sm.LoadAndDelete(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}

	if sm.Exists(1) {
		t.Errorf("expected key 1 to be deleted")
	}
}

func TestSafeMap_LoadOrStore_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	value, loaded := sm.LoadOrStore(1, "one")
	if loaded || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}

	value, loaded = sm.LoadOrStore(1, "two")
	if !loaded || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}
}

func TestSafeMap_Swap_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")

	previous, loaded := sm.Swap(1, "two")
	if !loaded || previous != "one" {
		t.Errorf("expected 'one', got '%v'", previous)
	}

	value, ok := sm.Load(1)
	if !ok || value != "two" {
		t.Errorf("expected 'two', got '%v'", value)
	}
}

func TestSafeMap_Range_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")
	sm.Store(2, "two")

	keys := make(map[int]bool)
	sm.Range(func(k int, v string) bool {
		keys[k] = true
		return true
	})

	if len(keys) != 2 || !keys[1] || !keys[2] {
		t.Errorf("expected keys 1 and 2, got '%v'", keys)
	}
}

func TestSafeMap_Clear_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")
	sm.Clear()

	if sm.Len() != 0 {
		t.Errorf("expected length 0, got '%v'", sm.Len())
	}
}

func TestSafeMap_Keys_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")
	sm.Store(2, "two")

	keys := sm.Keys()
	if len(keys) != 2 || (keys[0] != 1 && keys[1] != 2) {
		t.Errorf("expected keys 1 and 2, got '%v'", keys)
	}
}

func TestSafeMap_KeysInRange_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Store(1, "one")
	sm.Store(2, "two")
	sm.Store(3, "three")

	keys := sm.KeysInRange(1, 3)
	if len(keys) != 2 || keys[0] != 1 || keys[1] != 2 {
		t.Errorf("expected keys 1 and 2, got '%v'", keys)
	}
}

func BenchmarkSafeMap_Store(b *testing.B) {
	sm := New[int, string]()
	for i := 0; i < b.N; i++ {
		sm.Store(i, "value")
	}
}

func BenchmarkSafeMap_Load(b *testing.B) {
	sm := New[int, string]()
	sm.Store(1, "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Load(1)
	}
}

func BenchmarkSafeMap_Delete(b *testing.B) {
	sm := New[int, string]()
	for i := 0; i < b.N; i++ {
		sm.Store(i, "value")
		sm.Delete(i)
	}
}

func BenchmarkSafeMap_Ordered_Store(b *testing.B) {
	sm := New[int, string](Ordered)
	for i := 0; i < b.N; i++ {
		sm.Store(i, "value")
	}
}

func BenchmarkSafeMap_Ordered_Load(b *testing.B) {
	sm := New[int, string](Ordered)
	sm.Store(1, "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Load(1)
	}
}

func BenchmarkSafeMap_Ordered_Delete(b *testing.B) {
	sm := New[int, string](Ordered)
	for i := 0; i < b.N; i++ {
		sm.Store(i, "value")
		sm.Delete(i)
	}
}

func BenchmarkSyncMap_Store(b *testing.B) {
	var sm sync.Map
	for i := 0; i < b.N; i++ {
		sm.Store(i, "value")
	}
}

func BenchmarkSyncMap_Load(b *testing.B) {
	var sm sync.Map
	sm.Store(1, "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Load(1)
	}
}

func BenchmarkSyncMap_Delete(b *testing.B) {
	var sm sync.Map
	for i := 0; i < b.N; i++ {
		sm.Store(i, "value")
		sm.Delete(i)
	}
}

