package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileCacheInDir(t *testing.T) {
	tmpDir := t.TempDir()
	cacheDir := filepath.Join(tmpDir, "test-cache")

	c, err := NewFileCacheInDir(cacheDir)
	if err != nil {
		t.Fatalf("NewFileCacheInDir failed: %v", err)
	}

	if c.Dir() != cacheDir {
		t.Errorf("Expected dir %s, got %s", cacheDir, c.Dir())
	}

	info, err := os.Stat(cacheDir)
	if err != nil {
		t.Fatalf("Cache directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Cache path is not a directory")
	}
}

func TestCacheSetAndGet(t *testing.T) {
	c, err := NewFileCacheInDir(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	data := testData{Name: "test", Value: 42}

	if err := c.Set("mykey", "v1", data); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get with correct version
	raw, ok := c.Get("mykey", "v1")
	if !ok {
		t.Fatal("Expected cache hit")
	}

	var result testData
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if result.Name != "test" || result.Value != 42 {
		t.Errorf("Got %+v, want {Name:test Value:42}", result)
	}
}

func TestCacheMissOnVersionChange(t *testing.T) {
	c, err := NewFileCacheInDir(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Set("key", "commit-abc", "old-data"); err != nil {
		t.Fatal(err)
	}

	// Same key, different version → miss
	_, ok := c.Get("key", "commit-def")
	if ok {
		t.Error("Expected cache miss when version changed")
	}

	// Update with new version
	if err := c.Set("key", "commit-def", "new-data"); err != nil {
		t.Fatal(err)
	}

	raw, ok := c.Get("key", "commit-def")
	if !ok {
		t.Fatal("Expected cache hit after update")
	}
	var result string
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatal(err)
	}
	if result != "new-data" {
		t.Errorf("Got %q, want %q", result, "new-data")
	}
}

func TestCacheMissOnNonexistentKey(t *testing.T) {
	c, err := NewFileCacheInDir(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	_, ok := c.Get("nonexistent", "v1")
	if ok {
		t.Error("Expected cache miss for nonexistent key")
	}
}

func TestCacheMissOnCorruptedFile(t *testing.T) {
	dir := t.TempDir()
	c, err := NewFileCacheInDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Write valid entry first to get the path
	if err := c.Set("key", "v1", "data"); err != nil {
		t.Fatal(err)
	}

	// Corrupt the cache file
	path := c.pathFor("key")
	if err := os.WriteFile(path, []byte("not valid json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, ok := c.Get("key", "v1")
	if ok {
		t.Error("Expected cache miss for corrupted file")
	}
}

func TestCacheDifferentKeys(t *testing.T) {
	c, err := NewFileCacheInDir(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Set("key1", "v1", "data1"); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("key2", "v1", "data2"); err != nil {
		t.Fatal(err)
	}

	raw1, ok := c.Get("key1", "v1")
	if !ok {
		t.Fatal("Expected hit for key1")
	}
	raw2, ok := c.Get("key2", "v1")
	if !ok {
		t.Fatal("Expected hit for key2")
	}

	var d1, d2 string
	if err := json.Unmarshal(raw1, &d1); err != nil {
		t.Fatalf("Unmarshal raw1 failed: %v", err)
	}
	if err := json.Unmarshal(raw2, &d2); err != nil {
		t.Fatalf("Unmarshal raw2 failed: %v", err)
	}

	if d1 != "data1" || d2 != "data2" {
		t.Errorf("Got %q and %q, want data1 and data2", d1, d2)
	}
}

func TestCacheComplexData(t *testing.T) {
	c, err := NewFileCacheInDir(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	type LangInfo struct {
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		Percent int    `json:"percent"`
	}

	langs := []LangInfo{
		{Name: "Go", Size: 15000, Percent: 75},
		{Name: "Python", Size: 5000, Percent: 25},
	}

	if err := c.Set("/home/user/src/myproject", "abc123def", langs); err != nil {
		t.Fatal(err)
	}

	raw, ok := c.Get("/home/user/src/myproject", "abc123def")
	if !ok {
		t.Fatal("Expected cache hit")
	}

	var result []LangInfo
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(result))
	}
	if result[0].Name != "Go" || result[0].Size != 15000 {
		t.Errorf("Unexpected first entry: %+v", result[0])
	}
}
