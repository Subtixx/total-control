package plugins

func (p *PluginInfo) Can(capability string) bool {
	if p.Capabilities == nil {
		return false
	}

	for _, canCapability := range p.Capabilities {
		if canCapability == capability {
			return true
		}
	}
	return false
}

func (p *PluginInfo) CanAccessFileSystem() bool {
	return p.Can(CapabilityFileSystem)
}

func (p *PluginInfo) CanAccessNetwork() bool {
	return p.Can(CapabilityNetwork)
}
