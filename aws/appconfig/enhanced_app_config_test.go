package awsappconfig

import (
	"context"
	"fmt"
	awsappconfigadvance "github.com/hxy1991/sdk-go/aws/appconfig/advance"
	"github.com/hxy1991/sdk-go/constant"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// The 'app1' application with the 'test' environment must have been created in "us-east-1" region
var regionName = "us-east-1"
var applicationName = "app1"
var environmentName = "Test"

func TestAppConfig_GetConfiguration(t *testing.T) {
	type fields struct {
		applicationName string
		environmentName string
		regionName      string
	}
	type args struct {
		configurationName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "case1",
			fields: fields{
				regionName:      regionName,
				applicationName: applicationName,
				environmentName: environmentName,
			},
			args: args{
				configurationName: fmt.Sprintf("%d", time.Now().UnixNano()),
			},
			want:    time.Now().Format(time.RFC3339),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(t, constant.RegionEnvName, tt.fields.regionName)
			setEnv(t, constant.EnvironmentEnvName, tt.fields.environmentName)

			appConfig, appConfigAdvance, err := newConfig4Test(tt.fields.applicationName)
			assert.Nil(t, err)

			// 新增
			isSuccess, err := appConfigAdvance.CreateConfiguration(context.TODO(), tt.args.configurationName, tt.want)
			assert.True(t, (err != nil) == tt.wantErr, "CreateConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			assert.True(t, isSuccess, "CreateConfiguration() fail")

			time.Sleep(appConfig.cacheRefreshTicker.Interval() + appConfig.timeout)

			// 查询，从缓存中获取
			got, err := appConfig.GetConfiguration(context.TODO(), tt.args.configurationName)
			assert.True(t, (err != nil) == tt.wantErr, "GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "GetConfiguration() got = %v, want %v", got, tt.want)

			enhancedConfiguration, err := appConfig.GetEnhancedConfiguration(context.TODO(), tt.args.configurationName)
			assert.True(t, (err != nil) == tt.wantErr, "GetEnhancedConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			assert.True(t, enhancedConfiguration.IsCache)

			// 删除
			isSuccess, err = appConfigAdvance.DeleteConfiguration(context.TODO(), tt.args.configurationName)
			assert.True(t, (err != nil) == tt.wantErr, "DeleteConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			assert.True(t, isSuccess, "DeleteConfiguration() fail")
		})
	}
}

func TestAppConfig_GetConfigurationIgnoreCache(t *testing.T) {
	type fields struct {
		regionName      string
		applicationName string
		environmentName string
	}
	type args struct {
		configurationName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "case1",
			fields: fields{
				regionName:      regionName,
				applicationName: applicationName,
				environmentName: environmentName,
			},
			args: args{
				configurationName: fmt.Sprintf("%d", time.Now().UnixNano()),
			},
			want:    time.Now().Format(time.RFC3339),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(t, constant.RegionEnvName, tt.fields.regionName)
			setEnv(t, constant.EnvironmentEnvName, tt.fields.environmentName)

			appConfig, appConfigAdvance, err := newConfig4Test(tt.fields.applicationName)
			assert.Nil(t, err)

			t.Log("configurationName: ", tt.args.configurationName)
			// 新增
			isSuccess, err := appConfigAdvance.CreateConfiguration(context.TODO(), tt.args.configurationName, tt.want)
			assert.True(t, (err != nil) == tt.wantErr, "CreateConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			assert.True(t, isSuccess, "CreateConfiguration() fail")

			// 查询，不从缓存中获取
			got, err := appConfig.GetConfigurationIgnoreCache(context.TODO(), tt.args.configurationName)
			assert.True(t, (err != nil) == tt.wantErr, "GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "GetConfiguration() got = %v, want %v", got, tt.want)

			// 删除
			isSuccess, err = appConfigAdvance.DeleteConfiguration(context.TODO(), tt.args.configurationName)
			assert.True(t, (err != nil) == tt.wantErr, "DeleteConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			assert.True(t, isSuccess, "DeleteConfiguration() fail")
		})
	}
}

func TestAppConfig_UpdateIsCache(t *testing.T) {
	setEnvs(t)

	appConfig, err := NewWithOptions(WithApplicationName(applicationName))
	assert.Nil(t, err)

	configurationName := fmt.Sprintf("TestAppConfig_UpdateIsCache-%d", time.Now().Unix())

	createConfiguration(t, configurationName)

	// from aws app config
	getConfiguration(t, appConfig, configurationName, false)
	// from cache
	getConfiguration(t, appConfig, configurationName, true)

	err = appConfig.ApplyWithOptions(WithIsCache(false))
	assert.Nil(t, err)

	assert.Nil(t, appConfig.cache)
	assert.Nil(t, appConfig.cacheRefreshTicker)

	// from aws app config
	getConfiguration(t, appConfig, configurationName, false)
	getConfiguration(t, appConfig, configurationName, false)

	err = appConfig.ApplyWithOptions(WithIsCache(true))
	assert.Nil(t, err)

	assert.NotNil(t, appConfig.cache)
	assert.NotNil(t, appConfig.cacheRefreshTicker)

	// from aws app config
	getConfiguration(t, appConfig, configurationName, false)
	// from cache
	getConfiguration(t, appConfig, configurationName, true)

	deleteConfiguration(t, configurationName)
}

func TestAppConfig_UpdateCacheRefreshInterval(t *testing.T) {
	setEnvs(t)

	appConfig, err := NewWithApplicationName(applicationName)
	assert.Nil(t, err)

	configurationName := fmt.Sprintf("TestAppConfig_UpdateCacheRefreshInterval-%d", time.Now().Unix())

	createConfiguration(t, configurationName)

	// from aws app config
	getConfiguration(t, appConfig, configurationName, false)
	// from cache
	getConfiguration(t, appConfig, configurationName, true)

	newCacheRefreshInterval := time.Second * 10
	err = appConfig.ApplyWithOptions(WithCacheRefreshInterval(newCacheRefreshInterval))
	assert.Nil(t, err)

	assert.Equal(t, newCacheRefreshInterval, appConfig.cacheRefreshTicker.Interval(), "they should be equal")

	// from cache
	getConfiguration(t, appConfig, configurationName, true)

	deleteConfiguration(t, configurationName)
}

func TestAppConfig_UpdateCacheLimit(t *testing.T) {
	setEnvs(t)

	appConfig, err := NewWithApplicationName(applicationName)
	assert.Nil(t, err)

	configurationName1 := fmt.Sprintf("TestAppConfig_UpdateCacheLimit-1-%d", time.Now().Unix())

	createConfiguration(t, configurationName1)

	// from aws app config
	getConfiguration(t, appConfig, configurationName1, false)
	// from cache
	getConfiguration(t, appConfig, configurationName1, true)

	newCacheLimit := int64(1)
	err = appConfig.ApplyWithOptions(WithCacheLimit(newCacheLimit))
	assert.Nil(t, err)

	assert.Equal(t, newCacheLimit, appConfig.cache.CacheLimit(), "they should be equal")

	// from cache
	getConfiguration(t, appConfig, configurationName1, true)

	configurationName2 := fmt.Sprintf("TestAppConfig_UpdateCacheLimit-2-%d", time.Now().Unix())

	createConfiguration(t, configurationName2)

	getConfiguration(t, appConfig, configurationName2, false)
	getConfiguration(t, appConfig, configurationName2, true)

	getConfiguration(t, appConfig, configurationName1, false)
	getConfiguration(t, appConfig, configurationName1, true)

	deleteConfiguration(t, configurationName1)
	deleteConfiguration(t, configurationName2)
}

func TestAppConfig_UpdateTimeOut(t *testing.T) {
	setEnvs(t)

	appConfig, err := NewWithApplicationName(applicationName)
	assert.Nil(t, err)

	configurationName := fmt.Sprintf("TestAppConfig_UpdateTimeOut-%d", time.Now().Unix())

	createConfiguration(t, configurationName)

	// from aws app config
	getConfiguration(t, appConfig, configurationName, false)
	// from cache
	getConfiguration(t, appConfig, configurationName, true)

	newTimeOut := time.Second * 20
	err = appConfig.ApplyWithOptions(WithTimeout(newTimeOut))
	assert.Nil(t, err)

	assert.Equal(t, newTimeOut, appConfig.timeout, "they should be equal")

	// from cache
	getConfiguration(t, appConfig, configurationName, true)

	deleteConfiguration(t, configurationName)
}

func TestAppConfig_GetOptions(t *testing.T) {
	setEnvs(t)

	appConfig1, err := NewWithApplicationName(applicationName)
	assert.Nil(t, err)

	assert.NotNil(t, appConfig1)
	assert.NotNil(t, appConfig1.cache)
	assert.NotNil(t, appConfig1.cacheRefreshTicker)

	assert.Equal(t, appConfig1.regionName, regionName, "they should be equal")
	assert.Equal(t, appConfig1.applicationName, applicationName, "they should be equal")
	assert.Equal(t, appConfig1.environmentName, environmentName, "they should be equal")
	assert.Equal(t, appConfig1.isCache, defaultIsCache, "they should be equal")
	assert.Equal(t, appConfig1.cacheLimit, defaultCacheLimit, "they should be equal")
	assert.Equal(t, appConfig1.cacheRefreshInterval, defaultCacheRefreshInterval, "they should be equal")
	assert.Equal(t, appConfig1.timeout, defaultTimeout, "they should be equal")

	oregonRegionName := "us-west-2"

	appConfig2, err := NewWithOptions(
		WithApplicationName(applicationName),
		WithRegionName(oregonRegionName),
		WithIsCache(false),
	)
	assert.Nil(t, err)

	assert.NotNil(t, appConfig2)
	assert.Nil(t, appConfig2.cache)
	assert.Nil(t, appConfig2.cacheRefreshTicker)

	assert.Equal(t, appConfig2.regionName, oregonRegionName, "they should be equal")
	assert.Equal(t, appConfig2.applicationName, applicationName, "they should be equal")
	assert.Equal(t, appConfig2.environmentName, environmentName, "they should be equal")
	assert.Equal(t, appConfig2.isCache, false, "they should be equal")
	assert.Equal(t, appConfig2.cacheLimit, defaultCacheLimit, "they should be equal")
	assert.Equal(t, appConfig2.cacheRefreshInterval, defaultCacheRefreshInterval, "they should be equal")
	assert.Equal(t, appConfig2.timeout, defaultTimeout, "they should be equal")

	assert.Equal(t, appConfig1.regionName, regionName, "they should be equal")
	assert.Equal(t, appConfig1.applicationName, applicationName, "they should be equal")
	assert.Equal(t, appConfig1.environmentName, environmentName, "they should be equal")
	assert.Equal(t, appConfig1.isCache, defaultIsCache, "they should be equal")
	assert.Equal(t, appConfig1.cacheLimit, defaultCacheLimit, "they should be equal")
	assert.Equal(t, appConfig1.cacheRefreshInterval, defaultCacheRefreshInterval, "they should be equal")
	assert.Equal(t, appConfig1.timeout, defaultTimeout, "they should be equal")
}

func setEnvs(t *testing.T) {
	setEnv(t, constant.RegionEnvName, regionName)
	setEnv(t, constant.EnvironmentEnvName, environmentName)
}

func setEnv(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	assert.Nil(t, err)
}

func newConfig4Test(applicationName string) (*EnhancedAppConfig, *awsappconfigadvance.EnhancedAppConfigAdvance, error) {
	appConfig, err := NewWithOptions(
		WithCacheRefreshInterval(time.Second*10),
		WithApplicationName(applicationName),
	)
	if err != nil {
		return nil, nil, err
	}
	appConfigAdvance, err := awsappconfigadvance.NewWithOptions(
		awsappconfigadvance.WithApplicationName(applicationName),
	)
	return appConfig, appConfigAdvance, err
}

func createConfiguration(t *testing.T, configurationName string) {
	content := time.Now().Format(time.RFC3339)
	appConfigAdvance, err := awsappconfigadvance.NewWithOptions(
		awsappconfigadvance.WithApplicationName(applicationName),
	)
	assert.Nil(t, err)

	isSuccess, err := appConfigAdvance.CreateConfiguration(context.TODO(), configurationName, content)
	assert.Nil(t, err)

	assert.True(t, isSuccess, "createConfiguration fail")
}

func getConfiguration(t *testing.T, appConfig *EnhancedAppConfig, configurationName string, isFromCache bool) {
	configuration, err := appConfig.GetEnhancedConfiguration(context.TODO(), configurationName)
	assert.Nil(t, err)
	assert.Equal(t, isFromCache, configuration.IsCache, "expected %v, but received %v", isFromCache, configuration.IsCache)
}

func deleteConfiguration(t *testing.T, configurationName string) {
	appConfigAdvance, err := awsappconfigadvance.NewWithOptions(
		awsappconfigadvance.WithApplicationName(applicationName),
	)
	assert.Nil(t, err)

	isSuccess, err := appConfigAdvance.DeleteConfiguration(context.TODO(), configurationName)
	assert.Nil(t, err)

	assert.True(t, isSuccess, "deleteConfiguration fail")
}
