package markdown

import "github.com/nekrassov01/table/internal/scope"

// Scope selects the table parts to which an option applies. Values may be
// combined; zero and undefined bits select no parts.
//
// GFM tables have no footer section, so this package does not define
// ScopeFooter.
type Scope = scope.Scope

const (
	// ScopeHeader applies the option to the header row.
	ScopeHeader = scope.Header

	// ScopeBody applies the option to the data rows.
	ScopeBody = scope.Body
)
