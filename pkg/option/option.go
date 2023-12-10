package option

// Option holds generic options for Functional Options Pattern.
type Option[C any] func(C)
