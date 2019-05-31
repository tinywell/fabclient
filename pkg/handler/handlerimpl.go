package handler

// Impl implement for handler
type Impl struct {
}

// NewHandler return new Impl
func NewHandler() (*Impl, error) {
	return &Impl{}, nil
}
