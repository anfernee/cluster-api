package ipam

import (
	"reflect"
	"sort"
	"testing"
)

var testAddresses = []Address{
	{
		IP:       "1.1.1.1",
		Hostname: "1-1-1-1.local",
		DNS:      []string{"8.8.8.8"},
	},
	{
		IP:       "1.1.1.2",
		Hostname: "1-1-1-2.local",
		DNS:      []string{"8.8.8.8"},
	},
	{
		IP:       "1.1.1.3",
		Hostname: "1-1-1-3.local",
		DNS:      []string{"8.8.8.8"},
	},
}

func TestInMemStorageAllocate(t *testing.T) {
	s := NewInMemStorage()

	s.Add("cluster-1", testAddresses)

	var got []string
	for i := 0; i < len(testAddresses); i++ {
		a, err := s.Allocate("cluster-1")
		if err != nil {
			t.Fatalf("unexpected error from t.Allocate(cluster): %q", err)
		}
		got = append(got, a.IP)
	}

	sort.Strings(got)
	expected := []string{"1.1.1.1", "1.1.1.2", "1.1.1.3"}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected Allocate results: %s; got %s", expected, got)
	}
}

func TestInMemStorageRemove(t *testing.T) {
	s := NewInMemStorage()

	s.Add("cluster-1", testAddresses)

	for i := 0; i < len(testAddresses); i++ {
		_, err := s.Allocate("cluster-1")
		if err != nil {
			t.Fatalf("unexpected error from t.Allocate(cluster): %q", err)
		}
	}

	err := s.Remove("cluster-1")
	if err != nil {
		t.Fatalf("unexpected error from t.Remove(cluster): %q", err)
	}

	_, err = s.Allocate("cluster-1")
	if err != ErrKeyNotFound {
		t.Fatalf("expect error %q, got %q", ErrKeyNotFound, err)
	}
}

func TestInMemStorageOverAllocate(t *testing.T) {
	s := NewInMemStorage()

	s.Add("cluster-1", testAddresses)

	for i := 0; i < len(testAddresses); i++ {
		_, err := s.Allocate("cluster-1")
		if err != nil {
			t.Fatalf("unexpected error from t.Allocate(cluster): %q", err)
		}
	}

	_, err := s.Allocate("cluster-1")
	if err != ErrNotEnoughAddress {
		t.Fatalf("expect error %q, got %q", ErrNotEnoughAddress, err)
	}
}

func TestInMemStorageRelease(t *testing.T) {
	s := NewInMemStorage()

	s.Add("cluster-1", testAddresses)

	var addresses []Address
	for i := 0; i < len(testAddresses); i++ {
		a, err := s.Allocate("cluster-1")
		if err != nil {
			t.Fatalf("unexpected error from t.Allocate(cluster): %q", err)
		}
		addresses = append(addresses, *a)
	}

	for _, address := range addresses {
		err := s.Release("cluster-1", address.Key())
		if err != nil {
			t.Fatalf("unexpected error from t.Release(cluster, ip): %q", err)
		}
	}
}

func TestInMemStorageReleaseUnAllocated(t *testing.T) {
	s := NewInMemStorage()

	s.Add("cluster-1", testAddresses)

	var addresses []Address
	for i := 0; i < len(testAddresses); i++ {
		a, err := s.Allocate("cluster-1")
		if err != nil {
			t.Fatalf("unexpected error from t.Allocate(cluster): %q", err)
		}
		addresses = append(addresses, *a)
	}

	err := s.Release("cluster-1", "9.9.9.9")
	if err != ErrReleaseUnallocatedIP {
		t.Fatalf("expect error %q, got %q", ErrReleaseUnallocatedIP, err)
	}
}
