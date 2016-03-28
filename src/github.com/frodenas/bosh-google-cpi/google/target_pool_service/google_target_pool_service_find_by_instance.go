package targetpool

func (t GoogleTargetPoolService) FindByInstance(vmLink string, region string) (string, bool, error) {
	// Unfortunatelly, there is no direct way to find what target pool is attached to an instance,
	// so we need to list all target pools and look up for the instance
	targetPools, err := t.List(region)
	if err != nil {
		return "", false, err
	}

	for _, targetPool := range targetPools {
		for _, instance := range targetPool.Instances {
			if instance == vmLink {
				return targetPool.Name, true, nil
			}
		}
	}

	return "", false, nil
}
