package awsappconfig

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/xray"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/google/uuid"
	"github.com/hxy1991/sdk-go/cache"
	"github.com/hxy1991/sdk-go/constant"
	"github.com/hxy1991/sdk-go/log"
	"github.com/hxy1991/sdk-go/ticker"
)

const (
	defaultIsCache = true
	// the max number of configurations that can be cached
	defaultCacheLimit           = int64(500)
	defaultCacheRefreshInterval = time.Second * 300
	defaultTimeout              = time.Second * 10
)

type EnhancedAppConfig struct {
	applicationName string
	environmentName string
	clientId        string

	regionName           string
	isCache              bool          // 是否开启全局缓存
	cacheLimit           int64         // 最多缓存多少个配置
	cacheRefreshInterval time.Duration // 缓存刷新间隔
	timeout              time.Duration // 获取配置的超时时间

	isXRayEnable bool // 是否开启 X-Ray

	appConfigClient    *appconfig.AppConfig
	cache              *cache.Cache
	cacheRefreshTicker *ticker.Ticker
}

type EnhancedConfiguration struct {
	ClientConfigurationVersion *string
	Content                    *string
	IsCache                    bool
	CacheAt                    int64
}

func NewWithApplicationName(applicationName string) (*EnhancedAppConfig, error) {
	return NewWithOptions(WithApplicationName(applicationName))
}

func NewWithOptions(opts ...Option) (*EnhancedAppConfig, error) {
	appConfig := &EnhancedAppConfig{
		applicationName:      "",
		environmentName:      os.Getenv(constant.EnvironmentEnvName),
		clientId:             uuid.NewString(),
		regionName:           os.Getenv(constant.RegionEnvName),
		isCache:              defaultIsCache,
		cacheLimit:           defaultCacheLimit,
		cacheRefreshInterval: defaultCacheRefreshInterval,
		timeout:              defaultTimeout,
	}

	err := appConfig.ApplyWithOptions(opts...)
	if err != nil {
		return nil, err
	}

	if appConfig.regionName == "" {
		msg := fmt.Sprintf("missing required field: RegionName or set %s env", constant.RegionEnvName)
		return nil, errors.New(msg)
	}

	if appConfig.applicationName == "" {
		return nil, errors.New("missing required field: ApplicationName")
	}

	if appConfig.environmentName == "" {
		msg := fmt.Sprintf("missing required field: EnvironmentName or set %s env", constant.EnvironmentEnvName)
		return nil, errors.New(msg)
	}

	if appConfig.appConfigClient == nil {
		err = appConfig.initAppConfigClient()
		if err != nil {
			return nil, err
		}
	}

	if appConfig.cache == nil {
		if appConfig.isCache {
			appConfig.initCache()
		} else {
			log.Warn("cache is off, application name: ", appConfig.applicationName, ", environment name: ", appConfig.environmentName)
		}
	}

	return appConfig, nil
}

func (appConfig *EnhancedAppConfig) initCache() {
	log.Info("start init cache and ticker, cacheLimit: ", appConfig.cacheLimit, ", cacheRefreshInterval: ", appConfig.cacheRefreshInterval)
	appConfig.cache = cache.New(appConfig.cacheLimit)
	appConfig.initRefreshCacheTicker()
	log.Info("init cache and ticker end")
}

func (appConfig *EnhancedAppConfig) initAppConfigClient() error {
	awsConfig := aws.Config{
		Region: aws.String(appConfig.regionName),
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: awsConfig,
	}))

	appConfigClient := appconfig.New(sess)

	if appConfigClient == nil {
		return errors.New("can not init aws AppConfig client")
	}

	if appConfig.isXRayEnable {
		xray.AWS(appConfigClient.Client)
	}

	appConfig.appConfigClient = appConfigClient

	return nil
}

func (appConfig *EnhancedAppConfig) initRefreshCacheTicker() {
	cacheRefreshFunc := func() {
		var ctx context.Context
		if appConfig.isXRayEnable {
			_ctx, segment := xray.BeginSegment(context.Background(), "EnhancedAppConfig-CacheRefresh")
			defer segment.Close(nil)

			ctx = _ctx
		} else {
			ctx = context.Background()
		}

		startTime := time.Now()
		log.Debug("start refresh all the caches")
		var refreshCacheWaitGroup sync.WaitGroup
		for _, keyI := range appConfig.cache.Keys() {
			refreshCacheWaitGroup.Add(1)
			// 多协程并发获取
			appConfig.refreshKey(ctx, &refreshCacheWaitGroup, keyI)
		}
		refreshCacheWaitGroup.Wait()
		log.Debug("end refresh all the caches, cost: ", time.Since(startTime))
	}

	appConfig.cacheRefreshTicker = ticker.New(appConfig.cacheRefreshInterval, cacheRefreshFunc)
	appConfig.cacheRefreshTicker.Start()
}

func (appConfig *EnhancedAppConfig) refreshKey(ctx context.Context, refreshCacheWaitGroup *sync.WaitGroup, keyI interface{}) {
	go func() {
		defer func() {
			refreshCacheWaitGroup.Done()
			if e := recover(); e != nil {
				stack := string(debug.Stack())
				fmt.Println(stack)
				fmt.Println(e)
			}
		}()

		key := keyI.(string)
		appConfig.Refresh(ctx, key)
	}()

}

func (appConfig *EnhancedAppConfig) Refresh(ctx context.Context, key string) {
	log.Debug("start refresh cache [", key, "]")
	valueI, found := appConfig.cache.Get(key)
	if !found {
		return
	}
	if valueI == nil {
		log.Warn("refresh cache [", key, "] fail, valueI is nil, cache has been removed")
		return
	}

	clientConfigurationVersion := valueI.(*EnhancedConfiguration).ClientConfigurationVersion
	configuration, err := appConfig.getConfigurationWithVersion(ctx, key, clientConfigurationVersion)
	if err != nil {
		if strings.Contains(err.Error(), "could not be found for account") {
			log.Warn("refresh cache [", key, "] fail, configuration profile not exist, ", err)
			// 配置不存在了，删除缓存
			appConfig.cache.Delete(key)
			return
		}
		log.Error("refresh cache [", key, "] error ", err)
		return
	}

	if configuration == nil {
		msg := fmt.Sprintf("get from aws app config failed [%s]", key)
		log.Error(msg)
		return
	}

	if configuration.Content == nil {
		log.Debug("cache not change of configuration [", key, "]")
	} else {
		log.Warn("cache change of configuration [", key, "], new configuration version: ", *configuration.ClientConfigurationVersion)
		appConfig.cache.Add(key, configuration)
	}
	log.Debug("end refresh cache [", key, "]")
}

func (appConfig *EnhancedAppConfig) GetConfiguration(ctx context.Context, configurationName string) (string, error) {
	configuration, err := appConfig.GetEnhancedConfiguration(ctx, configurationName)
	if err != nil {
		return "", err
	}
	return *configuration.Content, nil
}

func (appConfig *EnhancedAppConfig) GetEnhancedConfiguration(ctx context.Context, configurationName string) (*EnhancedConfiguration, error) {
	// get from cache if cache is on
	if appConfig.cache != nil {
		cacheValue, found := appConfig.cache.Get(configurationName)
		if found {
			if cacheValue != nil {
				configuration := cacheValue.(*EnhancedConfiguration)
				return configuration, nil
			}
			log.Warn("get configuration from cache, but the value of cache is nil ", configurationName)
		}
	}

	configuration, err := appConfig.getConfigurationWithVersion(ctx, configurationName, nil)
	if err != nil {
		return nil, err
	}

	if configuration == nil || configuration.Content == nil {
		msg := fmt.Sprintf("get from aws app config failed [%s]", configurationName)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	// add to cache if cache is on
	if appConfig.cache != nil {
		log.Debug("add to cache ", configurationName)
		configuration.IsCache = true
		configuration.CacheAt = time.Now().Unix()
		appConfig.cache.Add(configurationName, configuration)
	}

	return &EnhancedConfiguration{
		ClientConfigurationVersion: configuration.ClientConfigurationVersion,
		Content:                    configuration.Content,
		IsCache:                    false,
	}, nil
}

func (appConfig *EnhancedAppConfig) GetConfigurationIgnoreCache(ctx context.Context, configurationName string) (string, error) {
	configuration, err := appConfig.getConfigurationWithVersion(ctx, configurationName, nil)
	if err != nil {
		return "", err
	}

	if configuration == nil || configuration.Content == nil {
		msg := fmt.Sprintf("get from aws app config failed [%s]", configurationName)
		log.Error(msg)
		return "", errors.New(msg)
	}

	return *(configuration.Content), err
}

func (appConfig *EnhancedAppConfig) GetEnhancedConfigurationIgnoreCache(ctx context.Context, configurationName string) (*EnhancedConfiguration, error) {
	configuration, err := appConfig.getConfigurationWithVersion(ctx, configurationName, nil)
	if err != nil {
		return nil, err
	}

	if configuration == nil || configuration.Content == nil {
		msg := fmt.Sprintf("get from aws app config failed [%s]", configurationName)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	return &EnhancedConfiguration{
		ClientConfigurationVersion: configuration.ClientConfigurationVersion,
		Content:                    configuration.Content,
		IsCache:                    false,
	}, nil
}

func (appConfig *EnhancedAppConfig) getConfigurationWithVersion(ctx context.Context, configurationName string, configurationVersion *string) (*EnhancedConfiguration, error) {
	configurationOutput, err := appConfig.getConfiguration(ctx, configurationName, configurationVersion)
	if err != nil {
		return nil, err
	}

	if len(configurationOutput.Content) == 0 {
		configuration := EnhancedConfiguration{
			ClientConfigurationVersion: configurationOutput.ConfigurationVersion,
			Content:                    nil,
		}
		return &configuration, nil
	}

	content := string(configurationOutput.Content)
	configuration := EnhancedConfiguration{
		ClientConfigurationVersion: configurationOutput.ConfigurationVersion,
		Content:                    &content,
	}
	return &configuration, nil
}

func (appConfig *EnhancedAppConfig) getConfiguration(ctx context.Context, configurationName string, configurationVersion *string) (*appconfig.GetConfigurationOutput, error) {
	input := appconfig.GetConfigurationInput{
		Application:   aws.String(appConfig.applicationName),
		Environment:   aws.String(appConfig.environmentName),
		ClientId:      aws.String(appConfig.clientId),
		Configuration: aws.String(configurationName),
	}
	if configurationVersion != nil {
		input.ClientConfigurationVersion = configurationVersion
	}
	now := time.Now()
	ctx, cancelFn := context.WithTimeout(ctx, appConfig.timeout)
	defer cancelFn()
	configuration, err := appConfig.appConfigClient.GetConfigurationWithContext(ctx, &input)
	//configuration, err := appConfigClient.GetConfiguration(&input)
	if err == nil {
		log.Debug("get configuration from aws app config successfully, name: ", configurationName, ", cost: ", time.Since(now))
	}
	return configuration, err
}

func (appConfig *EnhancedAppConfig) ApplyWithOptions(opts ...Option) error {
	for _, opt := range opts {
		err := opt.apply(appConfig)
		if err != nil {
			return err
		}
	}
	return nil
}
