//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

// Package codegen generates Web Component bundles for go-html slot maps.
//
// Use it at build time, or through the cmd/codegen CLI:
//
//	result := GenerateBundle(map[string]string{
//		"H": "site-header",
//		"C": "app-main",
//	})
//	bundle, _ := result.Value.(string)
package codegen
