package provider

// Provider defines methods of a provider.
type Provider interface {
	Init() error
	ScaleToZeroAfter(service string, timeout uint64) error
	ScaleToZero(service string) error
	Scale(service string, number uint64) error
}
