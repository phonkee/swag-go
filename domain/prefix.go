package domain

type PrefixOptions struct {
	// always provided description (other will be appended)
	Description string
}

// Defaults sets default values correctly (even fix invalid values)
func (p *PrefixOptions) Defaults() {
}
