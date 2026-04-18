// Package template provides placeholder rendering for vaultpipe.
//
// Placeholders follow the syntax:
//
//	{{vault:<secret-path>#<field>}}
//
// Example:
//
//	DSN=postgres://{{vault:secret/db#username}}:{{vault:secret/db#password}}@localhost/app
//
// The Renderer resolves placeholders against a pre-fetched secrets map,
// allowing composed values to be built before injecting them into a
// child process environment.
package template
