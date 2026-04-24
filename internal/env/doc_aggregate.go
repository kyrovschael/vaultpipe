// Package env — AggregateSnapshot
//
// AggregateSnapshot merges the values of matching keys across multiple
// [Snapshot] values, joining them with a configurable separator instead of
// letting later snapshots silently overwrite earlier ones.
//
// # Example
//
//	// Combine TAGS from three sources into a single comma-separated value.
//	 result, err := env.AggregateSnapshot(
//	     []env.Snapshot{
//	         {"TAGS": "alpha"},
//	         {"TAGS": "beta", "REGION": "us-east-1"},
//	         {"TAGS": "gamma"},
//	     },
//	     env.AggregateOptions{Separator: ",", SortValues: true},
//	 )
//	 // result["TAGS"]   == "alpha,beta,gamma"
//	 // result["REGION"] == "us-east-1"
package env
