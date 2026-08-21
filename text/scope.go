package text

import "github.com/nekrassov01/table/internal/scope"

// Scope selects the table parts to which an option applies. Values may be
// combined; zero and undefined bits select no parts.
type Scope = scope.Scope

const (
	// ScopeHeader applies the option to the header rows.
	ScopeHeader = scope.Header

	// ScopeBody applies the option to the body rows.
	ScopeBody = scope.Body

	// ScopeFooter applies the option to the footer rows.
	ScopeFooter = scope.Footer
)
