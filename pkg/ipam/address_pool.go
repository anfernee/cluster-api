package ipam

type addressPool struct {
	storage Storage
	key     string
}

func (p *addressPool) Allocate() (*Address, error) {
	return p.storage.Allocate(p.key)
}

func (p *addressPool) Release(address *Address) error {
	return p.storage.Release(p.key, address.Key())
}
