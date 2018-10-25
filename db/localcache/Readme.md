github: https://github.com/patrickmn/go-cache
doc: https://godoc.org/github.com/patrickmn/go-cache


// Add an item to the cache only if an item doesn't already exist for the given key,
// or if the existing item has expired. Returns an error otherwise.
func (c Cache) Add(k string, x interface{}, d time.Duration) error

// Get an item from the cache. Returns the item or nil, and a bool indicating whether the key was found.
func (c Cache) Get(k string) (interface{}, bool)

// Delete an item from the cache.
// Does nothing if the key is not in the cache.
func (c Cache) Delete(k string)

// Delete all expired items from the cache.
func (c Cache) DeleteExpired()

// Delete all items from the cache.
func (c Cache) Flush()

// GetWithExpiration returns an item and its expiration time from the cache.
// It returns the item or nil, the expiration time if one is set
// (if the item never expires a zero value for time.Time is returned),
// and a bool indicating whether the key was found.
func (c Cache) GetWithExpiration(k string) (interface{}, time.Time, bool)

// Returns the number of items in the cache.
// This may include items that have expired, but have not yet been cleaned up.
func (c Cache) ItemCount() int

// Copies all unexpired items in the cache into a new map and returns it.
func (c Cache) Items() map[string]Item

// Sets an (optional) function that is called with the key and value when an item is evicted from the cache.
//  (Including when it is deleted manually, but not when it is overwritten.) Set to nil to disable.
func (c Cache) OnEvicted(f func(string, interface{}))

// Set a new value for the cache key only if it already exists, and the existing item hasn't expired. Returns an error otherwise.
func (c Cache) Replace(k string, x interface{}, d time.Duration) error

// Add an item to the cache, replacing any existing item. If the duration is 0 (DefaultExpiration), the cache's default expiration time is used.

// If it is -1 (NoExpiration), the item never expires.
func (c Cache) Set(k string, x interface{}, d time.Duration)

// Add an item to the cache, replacing any existing item, using the default expiration.
func (c Cache) SetDefault(k string, x interface{})

func (c Cache) Save(w io.Writer) (err error)
func (c Cache) SaveFile(fname string) error
func (c Cache) Load(r io.Reader) error
func (c Cache) LoadFile(fname string) error