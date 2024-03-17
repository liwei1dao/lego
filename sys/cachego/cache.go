package cachego

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"
)

func (this *CacheGo) SetDefault(k string, x interface{}) {
	this.Set(k, x, this.options.Expiration)
}

func (this *CacheGo) Replace(k string, x interface{}, d time.Duration) error {
	this.mu.Lock()
	_, found := this.get(k)
	if !found {
		this.mu.Unlock()
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	this.set(k, x, d)
	this.mu.Unlock()
	return nil
}

func (this *CacheGo) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	this.mu.RLock()
	// "Inlining" of get and Expired
	item, found := this.items[k]
	if !found {
		this.mu.RUnlock()
		return nil, time.Time{}, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			this.mu.RUnlock()
			return nil, time.Time{}, false
		}

		// Return the item and the expiration time
		this.mu.RUnlock()
		return item.Object, time.Unix(0, item.Expiration), true
	}

	// If expiration <= 0 (i.e. no expiration time set) then return the item
	// and a zeroed time.Time
	this.mu.RUnlock()
	return item.Object, time.Time{}, true
}

func (this *CacheGo) Increment(k string, n int64) error {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return fmt.Errorf("Item %s not found", k)
	}
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) + int(n)
	case int8:
		v.Object = v.Object.(int8) + int8(n)
	case int16:
		v.Object = v.Object.(int16) + int16(n)
	case int32:
		v.Object = v.Object.(int32) + int32(n)
	case int64:
		v.Object = v.Object.(int64) + n
	case uint:
		v.Object = v.Object.(uint) + uint(n)
	case uintptr:
		v.Object = v.Object.(uintptr) + uintptr(n)
	case uint8:
		v.Object = v.Object.(uint8) + uint8(n)
	case uint16:
		v.Object = v.Object.(uint16) + uint16(n)
	case uint32:
		v.Object = v.Object.(uint32) + uint32(n)
	case uint64:
		v.Object = v.Object.(uint64) + uint64(n)
	case float32:
		v.Object = v.Object.(float32) + float32(n)
	case float64:
		v.Object = v.Object.(float64) + float64(n)
	default:
		this.mu.Unlock()
		return fmt.Errorf("The value for %s is not an integer", k)
	}
	this.items[k] = v
	this.mu.Unlock()
	return nil
}

func (this *CacheGo) IncrementFloat(k string, n float64) error {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return fmt.Errorf("Item %s not found", k)
	}
	switch v.Object.(type) {
	case float32:
		v.Object = v.Object.(float32) + float32(n)
	case float64:
		v.Object = v.Object.(float64) + n
	default:
		this.mu.Unlock()
		return fmt.Errorf("The value for %s does not have type float32 or float64", k)
	}
	this.items[k] = v
	this.mu.Unlock()
	return nil
}

func (this *CacheGo) IncrementInt(k string, n int) (int, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementInt8(k string, n int8) (int8, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int8)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int8", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementInt16(k string, n int16) (int16, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int16)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int16", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementInt32(k string, n int32) (int32, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int32)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int32", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementInt64(k string, n int64) (int64, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int64)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int64", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementUint(k string, n uint) (uint, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementUintptr(k string, n uintptr) (uintptr, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uintptr)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uintptr", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementUint8(k string, n uint8) (uint8, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint8)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint8", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementUint16(k string, n uint16) (uint16, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint16)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint16", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementUint32(k string, n uint32) (uint32, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint32)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint32", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementUint64(k string, n uint64) (uint64, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint64)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint64", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementFloat32(k string, n float32) (float32, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(float32)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an float32", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) IncrementFloat64(k string, n float64) (float64, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(float64)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an float64", k)
	}
	nv := rv + n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) Decrement(k string, n int64) error {
	// TODO: Implement Increment and Decrement more cleanly.
	// (Cannot do Increment(k, n*-1) for uints.)
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return fmt.Errorf("Item not found")
	}
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) - int(n)
	case int8:
		v.Object = v.Object.(int8) - int8(n)
	case int16:
		v.Object = v.Object.(int16) - int16(n)
	case int32:
		v.Object = v.Object.(int32) - int32(n)
	case int64:
		v.Object = v.Object.(int64) - n
	case uint:
		v.Object = v.Object.(uint) - uint(n)
	case uintptr:
		v.Object = v.Object.(uintptr) - uintptr(n)
	case uint8:
		v.Object = v.Object.(uint8) - uint8(n)
	case uint16:
		v.Object = v.Object.(uint16) - uint16(n)
	case uint32:
		v.Object = v.Object.(uint32) - uint32(n)
	case uint64:
		v.Object = v.Object.(uint64) - uint64(n)
	case float32:
		v.Object = v.Object.(float32) - float32(n)
	case float64:
		v.Object = v.Object.(float64) - float64(n)
	default:
		this.mu.Unlock()
		return fmt.Errorf("The value for %s is not an integer", k)
	}
	this.items[k] = v
	this.mu.Unlock()
	return nil
}

func (this *CacheGo) DecrementFloat(k string, n float64) error {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return fmt.Errorf("Item %s not found", k)
	}
	switch v.Object.(type) {
	case float32:
		v.Object = v.Object.(float32) - float32(n)
	case float64:
		v.Object = v.Object.(float64) - n
	default:
		this.mu.Unlock()
		return fmt.Errorf("The value for %s does not have type float32 or float64", k)
	}
	this.items[k] = v
	this.mu.Unlock()
	return nil
}

func (this *CacheGo) DecrementInt(k string, n int) (int, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) DecrementInt8(k string, n int8) (int8, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int8)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int8", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}
func (this *CacheGo) DecrementInt16(k string, n int16) (int16, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int16)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int16", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) DecrementInt32(k string, n int32) (int32, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int32)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int32", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) DecrementInt64(k string, n int64) (int64, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int64)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int64", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}
func (this *CacheGo) DecrementUint(k string, n uint) (uint, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) DecrementUintptr(k string, n uintptr) (uintptr, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uintptr)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uintptr", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) DecrementUint8(k string, n uint8) (uint8, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint8)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint8", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}
func (this *CacheGo) DecrementUint16(k string, n uint16) (uint16, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint16)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint16", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) DecrementUint32(k string, n uint32) (uint32, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint32)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint32", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}

func (this *CacheGo) DecrementUint64(k string, n uint64) (uint64, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint64)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint64", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}
func (this *CacheGo) DecrementFloat32(k string, n float32) (float32, error) {
	this.mu.Lock()
	v, found := this.items[k]
	if !found || v.Expired() {
		this.mu.Unlock()
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(float32)
	if !ok {
		this.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an float32", k)
	}
	nv := rv - n
	v.Object = nv
	this.items[k] = v
	this.mu.Unlock()
	return nv, nil
}
func (this *CacheGo) OnEvicted(f func(string, interface{})) {
	this.mu.Lock()
	this.onEvicted = f
	this.mu.Unlock()
}
func (this *CacheGo) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	this.mu.Lock()
	this.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
	this.mu.Unlock()
}
func (this *CacheGo) Add(k string, x interface{}, d time.Duration) error {
	this.mu.Lock()
	_, found := this.get(k)
	if found {
		this.mu.Unlock()
		return fmt.Errorf("Item %s already exists", k)
	}
	this.set(k, x, d)
	this.mu.Unlock()
	return nil
}
func (this *CacheGo) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	this.mu.Lock()
	for k, v := range this.items {
		// "Inlining" of expired
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := this.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	this.mu.Unlock()
	for _, v := range evictedItems {
		this.onEvicted(v.key, v.value)
	}
}

func (this *CacheGo) SaveFile(fname string) error {
	fp, err := os.Create(fname)
	if err != nil {
		return err
	}
	err = this.Save(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (this *CacheGo) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	this.mu.RLock()
	defer this.mu.RUnlock()
	for _, v := range this.items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&this.items)
	return
}

func (this *CacheGo) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	err := dec.Decode(&items)
	if err == nil {
		this.mu.Lock()
		defer this.mu.Unlock()
		for k, v := range items {
			ov, found := this.items[k]
			if !found || ov.Expired() {
				this.items[k] = v
			}
		}
	}
	return err
}

func (this *CacheGo) LoadFile(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	err = this.Load(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (this *CacheGo) Items() map[string]Item {
	this.mu.RLock()
	defer this.mu.RUnlock()
	m := make(map[string]Item, len(this.items))
	now := time.Now().UnixNano()
	for k, v := range this.items {
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		m[k] = v
	}
	return m
}

func (this *CacheGo) ItemCount() int {
	this.mu.RLock()
	n := len(this.items)
	this.mu.RUnlock()
	return n
}

func (this *CacheGo) Flush() {
	this.mu.Lock()
	this.items = map[string]Item{}
	this.mu.Unlock()
}

func (this *CacheGo) get(k string) (interface{}, bool) {
	item, found := this.items[k]
	if !found {
		return nil, false
	}
	// "Inlining" of Expired
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Object, true
}
func (this *CacheGo) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	this.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}
func (this *CacheGo) delete(k string) (interface{}, bool) {
	if this.onEvicted != nil {
		if v, found := this.items[k]; found {
			delete(this.items, k)
			return v.Object, true
		}
	}
	delete(this.items, k)
	return nil, false
}
