package ipam_test

import (
	"fmt"

	"k8s.io/cluster-api/pkg/ipam"
)

func ExampleIPManager() {
	storage := ipam.NewInMemStorage()
	im := ipam.NewIPManager(storage)
	im.Add("cluster-1", []ipam.Address{
		{
			IP:       "1.1.1.1",
			Hostname: "1-1-1-1.local",
			DNS:      []string{"8.8.8.8"},
		},
	})

	address, err := im.AddressPool("cluster-1").Allocate()
	fmt.Println(address, err)

	err = im.AddressPool("cluster-1").Release(address)
	fmt.Println(err)

	err = im.RemoveAll("cluster-1")
	fmt.Println(err)

	// Output: &{1.1.1.1 1-1-1-1.local [8.8.8.8]} <nil>
	// <nil>
	// <nil>
}
