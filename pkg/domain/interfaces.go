package domain

// Identifiers represents the identifying attributes of a taxonomy entity.
type Identifiers struct {
	Name string
	ID   string
}

// UnqSegKeys defines the interface for entities with unique identifiers.
type UnqSegKeys interface {
	GetIdentities() Identifiers
}
