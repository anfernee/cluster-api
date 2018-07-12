package ipam

import "errors"

type ipManager struct {
	storage Storage
}

func NewIPManager(storage Storage) IPManager {
	return &ipManager{storage}
}

func (m *ipManager) Add(clusterName string, addresses []Address) error {
	return m.storage.Add(clusterName, addresses)
}

func (m *ipManager) Remove(clusterName string, addresses []Address) error {
	return errors.New("not implemented yet")
}

func (m *ipManager) RemoveAll(clusterName string) error {
	return m.storage.Remove(clusterName)
}

func (m *ipManager) AddressPool(clusterName string) AddressPool {
	return &addressPool{m.storage, clusterName}
}
