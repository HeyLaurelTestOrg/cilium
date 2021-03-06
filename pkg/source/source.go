// SPDX-License-Identifier: Apache-2.0
// Copyright 2019 Authors of Cilium

package source

// Source describes the source of a definition
type Source string

const (
	// Unspec is used when the source is unspecified
	Unspec Source = "unspec"

	// Local is the source used for state derived from local agent state.
	// Local state has the strongest ownership and can only be overwritten
	// by other local state.
	Local Source = "local"

	// KVStore is the source used for state derived from a key value store.
	// State in the key value stored takes precedence over orchestration
	// system state such as Kubernetes.
	KVStore Source = "kvstore"

	// Kubernetes is the source used for state derived from Kubernetes
	Kubernetes Source = "k8s"

	// CustomResource is the source used for state derived from Kubernetes
	// custom resources
	CustomResource Source = "custom-resource"

	// Generated is the source used for generated state which can be
	// overwritten by all other sources
	Generated Source = "generated"
)

// AllowOverwrite returns true if new state from a particular source is allowed
// to overwrite existing state from another source
func AllowOverwrite(existing, new Source) bool {
	switch existing {

	// Local state can only be overwritten by other local state
	case Local:
		return new == Local

	// KVStore can be overwritten by other kvstore or local state
	case KVStore:
		return new == KVStore || new == Local

	// Custom-resource state can be overwritten by everything except
	// generated, unspecified and Kubernetes (non-CRD) state
	case CustomResource:
		return new != Generated && new != Unspec && new != Kubernetes

	// Kubernetes state can be overwritten by everything except generated
	// and unspecified state
	case Kubernetes:
		return new != Generated && new != Unspec

	// Generated can be overwritten by everything except by Unspecified
	case Generated:
		return new != Unspec

	// Unspecified state can be overwritten by everything
	case Unspec:
		return true
	}

	return true
}
