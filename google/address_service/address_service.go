package gaddress

type AddressService interface {
	Find(id string, region string) (Address, bool, error)
	FindByIP(ipAddress string) (Address, bool, error)
}
