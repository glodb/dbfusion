package queryoptions

// FindOptions provides options for performing database find operations.
type FindOptions struct {
	// ForceDB indicates whether to force the operation to use a specific database.
	ForceDB bool

	// CacheResult specifies whether to cache the result of the find operation.
	CacheResult bool
}
