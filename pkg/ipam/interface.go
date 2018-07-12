package ipam

import (
	"errors"
	"path/filepath"
)

var (
	// ErrKeyNotFound ..
	ErrKeyNotFound = errors.New("ipam: key not found")

	// ErrAddressNotFound
	ErrAddressNotFound = errors.New("ipam: address not found")

	// ErrEmptyAdressPool
	ErrEmptyAddressPool = errors.New("ipam: empty address pool")

	// ErrNotEnoughAddress
	ErrNotEnoughAddress = errors.New("ipam: no addresses can be allocated")

	// ErrReleaseUnallocatedIP
	ErrReleaseUnallocatedIP = errors.New("ipam: cannot release an ip that is not allocated by this pool")
)

// GlobalCluster is used when there is no address pool per cluster, so all
// clusters use addresses from this pool.
const GlobalCluster = ""

// IPManager manages IP addresses by reserving and allocating IP address in a per-cluster
// address pool.
type IPManager interface {
	// Add adds addresses into an address pool identified by cluster name.
	// If the address pool doesn't exist, create a new one.
	Add(clusterName string, addresses []Address) error

	// Remove removes address from an address pool identified by cluster name.
	// TODO: For IP managements after cluster is created. maybe not needed.
	Remove(clusterName string, addresses []Address) error

	// RemoveAll remove an address pool identified by cluster name.
	RemoveAll(clusterName string) error

	// AddressPool returns an address pool.
	AddressPool(clusterName string) AddressPool
}

// AddressPool represents an IP address pool to allocate IP/Hostname
// and to release it.
type AddressPool interface {
	// Allocate allocates an address from the pool
	Allocate() (*Address, error)

	// Release releases the address to the pool
	Release(address *Address) error
}

// Storage is the abstraction of how ipam persist the free/used addresses.
type Storage interface {
	// Add key
	Add(key string, addresses []Address) error

	// Remove key
	Remove(key string) error

	// Allocate allocates one address from
	Allocate(key string) (*Address, error)

	// Release
	Release(key string, addressKey string) error
}

// Snapshot is the snapshot of storage which contains all the records.
type Snapshot map[string][2][]Address

// ExtraStorage provide extra functionality on storage to take a snapshot
// of all records and recover from it. It is useful when you want to migrate
// between different storage. e.g. from local storage to an etcd cluster.
type ExtraStorage interface {
	Snapshot() (Snapshot, error)
	Recover(Snapshot) error

	Storage
}

// Address is the IP address and the related network configuration.
type Address struct {
	// IP is an IP address (IPv6 okay?)
	IP string

	// Hostname is the hostname
	// +optional
	Hostname string

	// DNS is the DNS servers for DNS lookup
	// +optional
	DNS []string
}

// Key returns address key to be saved
func (a *Address) Key() string {
	return a.IP
}

// TODO: move the following to other files

// PathKey the the function to generate a key for key/value store
func PathKey(base string, address *Address) string {
	return filepath.Join(base, address.Key())
}

// KeyFunc generates key from Address object
type KeyFunc func(base string, address *Address) string
