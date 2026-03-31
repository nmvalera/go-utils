package tag

import (
	"context"
)

type tagsKeyType string

// FromNamespaceContext returns the set of tags for the given namespace from the context.
func FromNamespaceContext(ctx context.Context, namespace string) Set {
	if set, ok := ctx.Value(tagsKeyType(namespace)).(Set); ok {
		return set
	}
	return EmptySet
}

// WithNamespaceTags returns a new context with the given added to the namespace set.
// The namespace set on the parent context is not modified.
func WithNamespaceTags(ctx context.Context, namespace string, tags ...*Tag) context.Context {
	set := FromNamespaceContext(ctx, namespace).WithTags(tags...)
	return context.WithValue(ctx, tagsKeyType(namespace), set)
}

// DefaultNamespace is the default namespace for tags
var DefaultNamespace = ""

// FromContext returns the set of tags from the context for the default namespace set
func FromContext(ctx context.Context) Set {
	return FromNamespaceContext(ctx, DefaultNamespace)
}

// WithTags returns a new context with the given tags added to the default namespace set.
func WithTags(ctx context.Context, tags ...*Tag) context.Context {
	return WithNamespaceTags(ctx, DefaultNamespace, tags...)
}

// WithComponent returns a new context with the component chained tag added to the default namespace set.
// If chained is true, the component tag will be chained to the existing component tag.
//
// Example:
// WithComponent(ctx, "component1", true) will return a new context with the component tag "component1"
// WithComponent(ctx, "component2", true) will return a new context with the component tag "component1.component2"
//
// If chained is false, the component tag will replace the existing component tag.
//
// Example:
// WithComponent(ctx, "component1", false) will return a new context with the component tag "component1"
// WithComponent(ctx, "component2", false) will return a new context with the component tag "component2"
func WithComponent(ctx context.Context, component string, chained bool) context.Context {
	return WithTags(ctx, Key("component").String(component).Chained(chained))
}
