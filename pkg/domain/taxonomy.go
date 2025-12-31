package domain

// Taxonomy is the root aggregate containing all taxonomy data.
type Taxonomy struct {
	ApiVersion string
	SegL1s     map[string]Seg
	SegsL2s    map[string]Seg
}
