package display

import (
	"net/url"
	"sort"
	"sync" // Note: AX-6 — sync.RWMutex for registry guard, no core wrapper in pinned core module
	"time"

	core "dappco.re/go"
)

const (
	maxStorageOriginBytes      = 512
	maxStorageBucketBytes      = 128
	maxStorageKeyBytes         = 1024
	maxStorageValueBytes       = 1 << 20
	maxStorageEntriesPerOrigin = 1024
	maxStorageBytesPerOrigin   = 16 << 20
	maxStorageSearchResults    = 200
)

type StorageEntry struct {
	Origin    string    `json:"origin"`
	Bucket    string    `json:"bucket"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StorageRegistry struct {
	mu      sync.RWMutex
	entries map[string]StorageEntry
	store   *storageStore
}

func NewStorageRegistry() *StorageRegistry {
	registry := &StorageRegistry{entries: make(map[string]StorageEntry)}
	registry.store = openStorageStore()
	registry.loadPersistedEntries()
	return registry
}

func openStorageStore() *storageStore {
	path := storageDatabasePath()
	if core.Trim(path) == "" {
		return nil
	}
	storeInstance, err := newStorageStore(path)
	if err != nil {
		core.Error(
			"storage registry init failed",
			"file_path", path,
			"step", "open",
			"err", core.E("display.storage.open", "failed to open storage store", err),
		)
		return nil
	}
	return storeInstance
}

func storageDatabasePath() string {
	if override := core.Trim(core.Env("CORE_GUI_STORAGE_PATH")); override != "" {
		return override
	}
	if core.Trim(core.Env("CORE_GUI_STORAGE_PERSIST")) == "" {
		return ":memory:"
	}
	home := core.Trim(core.Env("DIR_HOME"))
	if home == "" {
		return ":memory:"
	}
	return core.Path(home, ".core", "state", core.Concat("gui-storage-", core.Env("PID"), ".db"))
}

func storageOriginForPageURL(pageURL string) string {
	trimmed := core.Trim(pageURL)
	if trimmed == "" {
		return ""
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || core.Trim(parsed.Scheme) == "" {
		return ""
	}
	switch core.Lower(core.Trim(parsed.Scheme)) {
	case "http", "https":
		if parsed.Host == "" {
			return ""
		}
		return core.Concat(parsed.Scheme, "://", parsed.Host)
	case "core":
		if parsed.Host == "" {
			return "core://"
		}
		return core.Concat("core://", parsed.Host)
	case "file":
		if parsed.Path == "" {
			return ""
		}
		return core.Concat("file://", parsed.Path)
	default:
		if parsed.Host == "" {
			return ""
		}
		origin := core.Concat(parsed.Scheme, "://", parsed.Host)
		if parsed.Path != "" {
			origin = core.Concat(origin, parsed.Path)
		}
		origin = trimTrailingSlash(origin)
		if core.Trim(origin) == "://" {
			return trimmed
		}
		return origin
	}
}

func trimTrailingSlash(value string) string {
	for core.HasSuffix(value, "/") {
		value = core.TrimSuffix(value, "/")
	}
	return value
}

func makeStorageEntryKey(origin, bucket, key string) string {
	return storageCompositeKey(origin, bucket, key)
}

func storageCompositeKey(origin, bucket, key string) string {
	return core.JSONMarshalString([]string{origin, bucket, key})
}

func decodeStorageCompositeKey(value string) (string, string, string, bool) {
	var parts []string
	if result := core.JSONUnmarshalString(value, &parts); !result.OK || len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

func (r *StorageRegistry) loadPersistedEntries() {
	if r == nil || r.store == nil {
		return
	}
	if r.entries == nil {
		r.entries = make(map[string]StorageEntry)
	}
	items, err := r.store.getAll("storage")
	if err != nil {
		return
	}
	for key, value := range items {
		var stored StorageEntry
		if result := core.JSONUnmarshalString(value, &stored); !result.OK {
			if origin, bucket, key, ok := decodeStorageCompositeKey(key); ok {
				stored = StorageEntry{
					Origin:    origin,
					Bucket:    bucket,
					Key:       key,
					Value:     value,
					UpdatedAt: time.Now(),
				}
			} else {
				continue
			}
		}
		if stored.UpdatedAt.IsZero() {
			stored.UpdatedAt = time.Now()
		}
		r.entries[makeStorageEntryKey(stored.Origin, stored.Bucket, stored.Key)] = stored
	}
}

func (r *StorageRegistry) Set(origin, bucket, key, value string) bool {
	if r == nil {
		return false
	}
	if !validStorageField(origin, maxStorageOriginBytes) ||
		!validStorageField(bucket, maxStorageBucketBytes) ||
		!validStorageField(key, maxStorageKeyBytes) ||
		(len(value) > maxStorageValueBytes) {
		return false
	}
	origin = core.Trim(origin)
	bucket = core.Trim(bucket)
	key = core.Trim(key)
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.entries == nil {
		r.entries = make(map[string]StorageEntry)
	}
	composite := makeStorageEntryKey(origin, bucket, key)
	if !r.withinOriginQuotaLocked(origin, composite, StorageEntry{
		Origin: origin,
		Bucket: bucket,
		Key:    key,
		Value:  value,
	}) {
		return false
	}
	entry := StorageEntry{
		Origin:    origin,
		Bucket:    bucket,
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}
	if r.store != nil {
		if err := r.store.set("storage", storageCompositeKey(origin, bucket, key), core.JSONMarshalString(entry)); err != nil {
			return false
		}
	}
	r.entries[composite] = entry
	return true
}

func (r *StorageRegistry) Delete(origin, bucket, key string) bool {
	if r == nil {
		return false
	}
	if !validStorageField(origin, maxStorageOriginBytes) ||
		!validStorageField(bucket, maxStorageBucketBytes) ||
		!validStorageField(key, maxStorageKeyBytes) {
		return false
	}
	origin = core.Trim(origin)
	bucket = core.Trim(bucket)
	key = core.Trim(key)
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.entries == nil {
		r.entries = make(map[string]StorageEntry)
	}
	composite := makeStorageEntryKey(origin, bucket, key)
	if r.store != nil {
		if err := r.store.delete("storage", storageCompositeKey(origin, bucket, key)); err != nil {
			return false
		}
	}
	delete(r.entries, composite)
	return true
}

func (r *StorageRegistry) Get(origin, bucket, key string) (StorageEntry, bool) {
	if r == nil {
		return StorageEntry{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.entries == nil {
		return StorageEntry{}, false
	}

	if entry, ok := r.entries[makeStorageEntryKey(origin, bucket, key)]; ok {
		return entry, true
	}

	var latest StorageEntry
	found := false
	for _, entry := range r.entries {
		if bucket != "" && entry.Bucket != bucket {
			continue
		}
		if key != "" && entry.Key != key {
			continue
		}
		if origin != "" && entry.Origin != origin {
			continue
		}
		if !found || entry.UpdatedAt.After(latest.UpdatedAt) {
			latest = entry
			found = true
		}
	}
	return latest, found
}

func (r *StorageRegistry) Search(query string) []StorageEntry {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.entries == nil {
		return nil
	}
	needle := core.Lower(core.Trim(query))
	results := make([]StorageEntry, 0)
	for _, entry := range r.entries {
		if needle == "" ||
			core.Contains(core.Lower(entry.Origin), needle) ||
			core.Contains(core.Lower(entry.Bucket), needle) ||
			core.Contains(core.Lower(entry.Key), needle) ||
			core.Contains(core.Lower(entry.Value), needle) {
			results = append(results, entry)
			if len(results) >= maxStorageSearchResults {
				break
			}
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})
	return results
}

func validStorageField(value string, limit int) bool {
	trimmed := core.Trim(value)
	return trimmed != "" && len(trimmed) <= limit
}

func (r *StorageRegistry) Snapshot(pageURL string) map[string]map[string]string {
	if r == nil {
		return map[string]map[string]string{}
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.entries == nil {
		return map[string]map[string]string{}
	}

	origin := storageOriginForPageURL(pageURL)
	if core.Trim(origin) == "" {
		return map[string]map[string]string{}
	}
	snapshot := make(map[string]map[string]string)
	for _, entry := range r.entries {
		if origin != "" && !storageEqualFold(entry.Origin, origin) {
			continue
		}
		bucket := snapshot[entry.Bucket]
		if bucket == nil {
			bucket = make(map[string]string)
			snapshot[entry.Bucket] = bucket
		}
		bucket[entry.Key] = entry.Value
	}
	return snapshot
}

func (r *StorageRegistry) withinOriginQuotaLocked(origin, ignoreComposite string, candidate StorageEntry) bool {
	if r == nil || r.entries == nil {
		return true
	}
	entries := 0
	bytes := 0
	for composite, entry := range r.entries {
		if !storageEqualFold(entry.Origin, origin) {
			continue
		}
		if composite == ignoreComposite {
			continue
		}
		entries++
		bytes += storageEntrySizeBytes(entry)
	}
	entries++
	bytes += storageEntrySizeBytes(candidate)
	if entries > maxStorageEntriesPerOrigin {
		return false
	}
	if bytes > maxStorageBytesPerOrigin {
		return false
	}
	return true
}

func storageEntrySizeBytes(entry StorageEntry) int {
	return len(entry.Origin) + len(entry.Bucket) + len(entry.Key) + len(entry.Value)
}

func storageEqualFold(left, right string) bool {
	return core.Lower(left) == core.Lower(right)
}

func (r *StorageRegistry) Close() resultFailure {
	if r == nil || r.store == nil {
		return nil
	}
	return r.store.close()
}
