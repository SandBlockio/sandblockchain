package types

import "strings"

// Query endpoints supported by the surprise querier
const (
	QueryListBrandedTokens = "list"
	QueryGetBrandedToken = "get"
	QueryGetTotalSupply = "supply"
)

type QueryResFetch []string

// implement fmt.Stringer
func (n QueryResFetch) String() string {
	return strings.Join(n[:], "\n")
}