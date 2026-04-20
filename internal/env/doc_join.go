// Package env – JoinSnapshot
//
// JoinSnapshot merges multiple [Snapshot] values into a single snapshot.
//
// When the same key appears in more than one source snapshot its values are
// concatenated using a configurable separator (default: ","). This is useful
// for environment variables that are list-like, such as PATH or
// JAVA_TOOL_OPTIONS, where contributions from multiple sources should all be
// preserved rather than one silently overwriting another.
//
// Example:
//
//	base := env.Snapshot{{Key: "PATH", Value: "/usr/bin"}}
//	extra := env.Snapshot{{Key: "PATH", Value: "/opt/bin"}}
//
//	result := env.JoinSnapshot(
//		[]env.Snapshot{base, extra},
//		env.DefaultJoinOptions(),
//	)
//	// result[0] == {Key:"PATH", Value:"/usr/bin,/opt/bin"}
package env
