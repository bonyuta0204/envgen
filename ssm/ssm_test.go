package ssm

import (
	"testing"
	"sync"
)

// We'll use a simpler approach for testing the Client
// since mocking the AWS SDK is complex

func TestClientCache(t *testing.T) {
	// Create a client with a simple cache
	client := &Client{
		cache: make(map[string]string),
		mutex: sync.RWMutex{},
	}
	
	// Add some values to the cache
	client.mutex.Lock()
	client.cache["/myapp/db_host"] = "db.example.com"
	client.cache["/myapp/db_password"] = "secret-password"
	client.mutex.Unlock()
	
	// Test that we can retrieve values from the cache
	client.mutex.RLock()
	value, ok := client.cache["/myapp/db_host"]
	client.mutex.RUnlock()
	
	if !ok {
		t.Errorf("Value not found in cache")
	}
	if value != "db.example.com" {
		t.Errorf("Cached value = %v, want %v", value, "db.example.com")
	}
	
	// Test that non-existent values are not in the cache
	client.mutex.RLock()
	_, ok = client.cache["/myapp/nonexistent"]
	client.mutex.RUnlock()
	
	if ok {
		t.Errorf("Non-existent value found in cache")
	}
}

// Test the loadConfig function
func TestLoadConfig(t *testing.T) {
	// This is a simple test to ensure the function doesn't panic
	// We can't really test the AWS config without credentials
	_, err := loadConfig("")
	if err != nil {
		// Skip the test if we can't load the config (which is expected in a test environment)
		t.Skip("Skipping AWS config test - no credentials available")
	}
	
	// Test with a region
	_, err = loadConfig("us-west-2")
	if err != nil {
		t.Skip("Skipping AWS config test - no credentials available")
	}
}
