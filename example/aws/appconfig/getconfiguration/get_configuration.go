package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/xray"
	awsappconfig "github.com/hxy1991/sdk-go/aws/appconfig"
	awsappconfigadvance "github.com/hxy1991/sdk-go/aws/appconfig/advance"
	"github.com/hxy1991/sdk-go/constant"
	"github.com/hxy1991/sdk-go/log"
	"os"
	"time"
)

var regionName = "us-east-1"
var applicationName = "app1"
var environmentName = "Test"

func main() {
	if setEnv(constant.RegionEnvName, regionName) {
		return
	}

	if setEnv(constant.EnvironmentEnvName, environmentName) {
		return
	}

	ctx, segment := xray.BeginSegment(context.Background(), "Example-GetConfiguration")
	defer segment.Close(nil)

	configurationName := fmt.Sprintf("get-configuration-%d", time.Now().Unix())
	createConfiguration(ctx, configurationName)
	getConfiguration(ctx, configurationName)
	deleteConfiguration(ctx, configurationName)
}

func createConfiguration(ctx context.Context, configurationName string) {
	appConfigAdvance, err := awsappconfigadvance.NewWithApplicationName(applicationName)
	if err != nil {
		panic(err)
	}

	ok, err := appConfigAdvance.CreateConfiguration(ctx, configurationName, time.Now().Format(time.RFC3339))
	if err != nil {
		panic(err)
	}

	log.Info("CreateConfiguration ", configurationName, " ", ok)
}

func getConfiguration(ctx context.Context, configurationName string) {
	appConfig, err := awsappconfig.NewWithOptions(
		awsappconfig.WithApplicationName(applicationName),
		awsappconfig.WithCacheRefreshInterval(time.Second*30),
		awsappconfig.WithRegionName(regionName),
	)
	if err != nil {
		panic(err)
	}

	content, err := appConfig.GetConfiguration(ctx, configurationName)
	if err != nil {
		log.Error(err)
	} else {
		log.Info(configurationName, ": ", content)
	}
}

func deleteConfiguration(ctx context.Context, configurationName string) {
	appConfigAdvance, err := awsappconfigadvance.NewWithApplicationName(applicationName)
	if err != nil {
		panic(err)
	}

	ok, err := appConfigAdvance.DeleteConfiguration(ctx, configurationName)
	if err != nil {
		panic(err)
	}

	log.Info("DeleteConfiguration ", configurationName, " ", ok)
}

func setEnv(key, value string) bool {
	err := os.Setenv(key, value)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}
