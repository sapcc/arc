package types

import (
	"fmt"
)

// Endpoints is a custom flag Var representing a list of transport endpoints
type Endpoints []string

// String returns the string representation of a endpoints var.
func (n *Endpoints) String() string {
	return fmt.Sprintf("%s", *n)
}

// Set appends the endpoint to the endpoints list.
func (n *Endpoints) Set(endpoint string) error {
	*n = append(*n, endpoint)
	return nil
}
