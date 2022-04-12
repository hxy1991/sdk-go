package awsappconfigadvance

type Option interface {
	apply(*EnhancedAppConfigAdvance) error
}

type optionFunc func(*EnhancedAppConfigAdvance) error

func (f optionFunc) apply(appConfigAdvance *EnhancedAppConfigAdvance) error {
	return f(appConfigAdvance)
}

func WithRegionName(regionName string) Option {
	return optionFunc(func(appConfigAdvance *EnhancedAppConfigAdvance) error {
		appConfigAdvance.regionName = regionName

		err := appConfigAdvance.initAppConfigClient()
		if err != nil {
			return err
		}

		return nil
	})
}

func WithApplicationName(applicationName string) Option {
	return optionFunc(func(appConfigAdvance *EnhancedAppConfigAdvance) error {
		appConfigAdvance.applicationName = applicationName
		return nil
	})
}

func WithEnvironmentName(environmentName string) Option {
	return optionFunc(func(appConfigAdvance *EnhancedAppConfigAdvance) error {
		appConfigAdvance.environmentName = environmentName
		return nil
	})
}
