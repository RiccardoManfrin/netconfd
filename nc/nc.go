package nc

import (
	"github.com/milosgajdos/tenus"
	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/logger"
)

//Create generic link interface
func CreateLink(kind string, ifname string, ifindex int) error {
	switch kind {
	case "dummy":
		{
			return CreateDummyLink(ifname, ifindex)
		}
	default:
		logger.Log.Fatal("Unknown Link Type " + kind)
	}
	return nil
}

//CreateDummyLink Creates a new dummy link
func CreateDummyLink(ifname string, ifindex int) error {
	dl, err := tenus.NewLink(ifname)
	if err != nil {
		logger.Log.Fatal(err)
	}
	if err = dl.SetLinkUp(); err != nil {
		logger.Log.Fatal(err)
	}
	return nil
}
