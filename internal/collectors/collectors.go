package collectors

type Collector interface {
	// Start begins the collection process.
	Start() error
	// Stop halts the collection process.
	Stop() error
}
