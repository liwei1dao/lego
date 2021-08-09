package registry

func newSys(options Options) (sys ISys, err error) {
	if options.RegistryType == Registry_Consul {
		sys, err = newConsul(options)
	} else if options.RegistryType == Registry_Nacos {
		sys, err = newNacos(options)
	}
	return
}
