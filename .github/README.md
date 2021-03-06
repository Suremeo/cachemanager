## How to install

Run `go get github.com/suremeo/cachemanager`

## Examples

```go
package main

import (
	"github.com/suremeo/cachemanager"
    "time"
)

var cache = cachemanager.NewCache().Run()

func main() {
	// Add an item, item can be anything, a structure, bytes, string, int, etc

	cache.Set("Suremeo", "cool guy")

	// Set how long items stay in cache (seconds) before they get automatically removed (preset is 30)

	cache.Expire = 60 * time.Second
    
    // Set how often to tick (Loops through all cached items and removes expired ones)
    
    cache.Tick = 1 * time.Second

	// fetch from cache

	item, _ := cache.Get("Suremeo")
	println(item.Data.(string))

	// Manually remove something from cache (usually it does it automatically after it expires)

	cache.Remove("Suremeo")

	// Cache file reading (if the file isn't already in the cache it read the file and adds it)

	cache.File("Booger.png")
	
	// Clear entire cache
	
	cache.Clear()
}
```
