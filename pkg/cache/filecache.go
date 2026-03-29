package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileCache provides a file-based cache using os.UserCacheDir for OS-agnostic storage.
// Each entry is a JSON file keyed by a hash of the cache key.
type FileCache struct {
	dir string
	mu  sync.Mutex
}

// CacheEntry wraps a cached value with a version key for invalidation.
type CacheEntry struct {
	Version string          `json:"version"`
	Data    json.RawMessage `json:"data"`
}

// NewFileCache creates a new file cache in the given subdirectory under os.UserCacheDir.
// For example, NewFileCache("allbctl", "languages") uses ~/.cache/allbctl/languages/ on Linux.
func NewFileCache(subDirs ...string) (*FileCache, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine cache directory: %w", err)
	}

	parts := append([]string{base}, subDirs...)
	dir := filepath.Join(parts...)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create cache directory %s: %w", dir, err)
	}

	return &FileCache{dir: dir}, nil
}

// NewFileCacheInDir creates a file cache in an explicit directory.
// Useful for testing.
func NewFileCacheInDir(dir string) (*FileCache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create cache directory %s: %w", dir, err)
	}
	return &FileCache{dir: dir}, nil
}

// Get retrieves a cached value. Returns the data and true if a valid entry
// exists with a matching version; otherwise returns nil and false.
func (c *FileCache) Get(key, version string) (json.RawMessage, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	path := c.pathFor(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	if entry.Version != version {
		return nil, false
	}

	return entry.Data, true
}

// Set stores a value in the cache with the given version key.
func (c *FileCache) Set(key, version string, data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal cache data: %w", err)
	}

	entry := CacheEntry{
		Version: version,
		Data:    jsonData,
	}

	entryData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("cannot marshal cache entry: %w", err)
	}

	path := c.pathFor(key)
	return os.WriteFile(path, entryData, 0644)
}

// pathFor returns the cache file path for a given key.
func (c *FileCache) pathFor(key string) string {
	hash := sha256.Sum256([]byte(key))
	filename := fmt.Sprintf("%x.json", hash[:16]) // 32 hex chars
	return filepath.Join(c.dir, filename)
}

// Dir returns the cache directory path.
func (c *FileCache) Dir() string {
	return c.dir
}
