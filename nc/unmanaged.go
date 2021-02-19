package nc

import "fmt"

type UnmanagedID string

// Unmanaged Unmanaged Resource For link type resources, the related context (routes or DHCP) are also unmanaged
type Unmanaged struct {
	// Resource type
	Type Type `json:"type,omitempty"`
	// ID of the resource to ignore
	ID UnmanagedID `json:"id,omitempty"`
}

var unmanaged = make(map[UnmanagedID]Unmanaged)

// Type identifies the type of a resource
type Type string

// List of scope
const (
	LINKTYPE Type = "link"
	DNSTYPE  Type = "dns"
)

//Print implements unmanaged resource print
func (u *Unmanaged) Print() string {
	return fmt.Sprintf("%v", u)
}

//UnamanagedListConfigure configures the array of unmanaged resources
func UnamanagedListConfigure(umgmts []Unmanaged) error {
	for _, u := range umgmts {
		err := UnmanagedDelete(u.ID)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok != true {
				return err
			}
		}
		if err := UnmanagedCreate(u); err != nil {
			return err
		}
	}
	return nil
}

//UnmanagedListDelete deletes all unmanaged resouces
func UnmanagedListDelete() error {
	umgmts, err := UnmanagedListGet()
	if err != nil {
		return err
	}
	for _, u := range umgmts {
		err = UnmanagedDelete(u.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

//UnmanagedListGet returns list of unmanaged objects
func UnmanagedListGet() ([]Unmanaged, error) {
	var u = make([]Unmanaged, len(unmanaged))
	i := 0
	for _, v := range unmanaged {
		u[i] = v
		i++
	}
	return u, nil
}

// UnmanagedCreate adds a new unmanaged network resource object
func UnmanagedCreate(u Unmanaged) error {
	switch u.Type {
	case LINKTYPE:
		{
			lid := LinkID(u.ID)
			_, err := LinkGet(lid)
			if err != nil {
				if _, ok := err.(*NotFoundError); ok == false {
					return mapNetlinkError(err, &u)
				}
			}
			unmanaged[u.ID] = u
		}
	case DNSTYPE:
		{
			if DnsID(u.ID) != DnsPrimary && DnsID(u.ID) != DnsSecondary {
				return NewUnknownUnsupportedDNSServersIDsError(DnsID(u.ID))
			}
			unmanaged[u.ID] = u
		}
	default:
		{
			return NewInvalidUnmanagedResourceTypeError(u.Type)
		}
	}
	return nil
}

// UnmanagedDelete adds a new unmanaged network resource object
func UnmanagedDelete(id UnmanagedID) error {
	if _, ok := unmanaged[id]; ok != true {
		return NewUnmanagedResourceNotFoundError(id)
	}
	delete(unmanaged, id)
	return nil
}

// UnmanagedGet adds a new unmanaged network resource object
func UnmanagedGet(id UnmanagedID) (Unmanaged, error) {
	u, ok := unmanaged[id]
	if ok != true {
		return Unmanaged{}, NewUnmanagedResourceNotFoundError(id)
	}
	return u, nil
}

func isUnmanaged(id UnmanagedID, t Type) bool {
	_, ok := unmanaged[id]
	return ok
}
