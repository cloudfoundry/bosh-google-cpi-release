package instancegroup

func (i GoogleInstanceGroupService) FindByInstance(vmLink string, zone string) (string, bool, error) {
	// Unfortunatelly, there is no direct way to find what instance group is attached to an instance,
	// so we need to list all instance groups and look up for the instance
	instanceGroups, err := i.List(zone)
	if err != nil {
		return "", false, err
	}

	for _, instanceGroup := range instanceGroups {
		for _, instance := range instanceGroup.Instances {
			if instance == vmLink {
				return instanceGroup.Name, true, nil
			}
		}
	}

	return "", false, nil
}
