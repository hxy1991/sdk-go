package awsappconfigadvance

import (
	"context"
	"errors"
	"fmt"
	"github.com/hxy1991/sdk-go/constant"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/aws/aws-xray-sdk-go/xray"
)

// AllAtOnceNotBake deployment strategy must have been created
const deploymentStrategyName = "AllAtOnceNotBake"

type EnhancedAppConfigAdvance struct {
	regionName      string
	applicationName string
	environmentName string
	applicationId   string
	environmentId   string
	appConfigClient *appconfig.AppConfig

	isXRayEnable bool
}

var applicationNameId = map[string]string{}
var environmentNameId = map[string]string{}

var configurationProfileNameId = sync.Map{}
var deploymentStrategyNameId = sync.Map{}

func NewWithApplicationName(applicationName string) (*EnhancedAppConfigAdvance, error) {
	return NewWithOptions(WithApplicationName(applicationName))
}

func NewWithOptions(opts ...Option) (*EnhancedAppConfigAdvance, error) {
	appConfigAdvance := &EnhancedAppConfigAdvance{
		regionName:      os.Getenv(constant.RegionEnvName),
		applicationName: "",
		environmentName: os.Getenv(constant.EnvironmentEnvName),

		applicationId:   "",
		environmentId:   "",
		appConfigClient: nil,
	}

	err := appConfigAdvance.ApplyOptions(opts...)
	if err != nil {
		return nil, err
	}

	if appConfigAdvance.applicationName == "" {
		return nil, errors.New("missing required field: ApplicationName")
	}

	if appConfigAdvance.regionName == "" {
		msg := fmt.Sprintf("missing required field: RegionName or set %s env", constant.RegionEnvName)
		return nil, errors.New(msg)
	}

	if appConfigAdvance.environmentName == "" {
		msg := fmt.Sprintf("missing required field: EnvironmentName or set %s env", constant.EnvironmentEnvName)
		return nil, errors.New(msg)
	}

	if appConfigAdvance.appConfigClient == nil {
		err = appConfigAdvance.initAppConfigClient()
		if err != nil {
			return nil, err
		}
	}

	var ctx context.Context
	if appConfigAdvance.isXRayEnable {
		_ctx, segment := xray.BeginSegment(context.Background(), "EnhancedAppConfigAdvance-NewWithOptions")
		defer segment.Close(nil)

		ctx = _ctx
	} else {
		ctx = context.Background()
	}

	err = appConfigAdvance.listApplications(ctx)
	if err != nil {
		return nil, err
	}

	err = appConfigAdvance.listDeploymentStrategies(ctx)
	if err != nil {
		return nil, err
	}

	err = appConfigAdvance.nameToId(ctx)
	if err != nil {
		return nil, err
	}

	err = appConfigAdvance.listConfigurationProfiles(ctx)
	if err != nil {
		return nil, err
	}

	return appConfigAdvance, err
}

func (appConfigAdvance *EnhancedAppConfigAdvance) initAppConfigClient() error {
	awsConfig := aws.Config{
		Region: aws.String(appConfigAdvance.regionName),
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: awsConfig,
	}))

	appConfigClient := appconfig.New(sess)

	if appConfigClient == nil {
		return errors.New("can not init aws AppConfig client")
	}

	if appConfigAdvance.isXRayEnable {
		xray.AWS(appConfigClient.Client)
	}

	appConfigAdvance.appConfigClient = appConfigClient
	return nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) UpdateConfiguration(ctx context.Context, configurationName string, content string) (bool, error) {
	configurationProfileId, found, err := appConfigAdvance.getConfigurationProfileId(ctx, configurationName)
	if err != nil {
		return false, err
	}
	if !found {
		msg := fmt.Sprintf("configuration [%s] do not exist in [%s] environment of [%s] application", configurationName, appConfigAdvance.environmentName, appConfigAdvance.applicationName)
		return false, errors.New(msg)
	}

	// 创建版本
	createHostedConfigurationVersionOutput, err := appConfigAdvance.createHostedConfigurationVersion(ctx, configurationProfileId, content)
	if err != nil {
		return false, err
	}

	// 发布版本
	configurationVersion := fmt.Sprintf("%d", *createHostedConfigurationVersionOutput.VersionNumber)
	startDeploymentOutput, err := appConfigAdvance.startDeployment(ctx, configurationProfileId, configurationVersion)
	if err != nil {
		return false, err
	}
	return startDeploymentOutput != nil, nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) getConfigurationProfileId(ctx context.Context, configurationName string) (string, bool, error) {
	configurationProfileId, found := configurationProfileNameId.Load(configurationName)
	if found {
		return configurationProfileId.(string), found, nil
	}
	// 再获取一次配置名称和ID的对应关系
	err := appConfigAdvance.listConfigurationProfiles(ctx)
	if err != nil {
		return "", false, err
	}

	configurationProfileId, found = configurationProfileNameId.Load(configurationName)
	if found {
		return configurationProfileId.(string), found, nil
	}

	return "", false, nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) CreateConfiguration(ctx context.Context, configurationName string, content string) (bool, error) {
	_, found, err := appConfigAdvance.getConfigurationProfileId(ctx, configurationName)
	if err != nil {
		return false, err
	}
	if found {
		// 配置已经存在
		msg := fmt.Sprintf("configuration [%s] already exist in [%s] environment of [%s] application", configurationName, appConfigAdvance.environmentName, appConfigAdvance.applicationName)
		return false, errors.New(msg)
	}

	// 创建配置 Profile
	createConfigurationProfileOutput, err := appConfigAdvance.createConfigurationProfile(ctx, configurationName)
	if err != nil {
		return false, err
	}
	configurationProfileId := *createConfigurationProfileOutput.Id

	// 创建版本
	createHostedConfigurationVersionOutput, err := appConfigAdvance.createHostedConfigurationVersion(ctx, configurationProfileId, content)
	if err != nil {
		return false, err
	}

	// 发布版本
	configurationVersion := fmt.Sprintf("%d", *createHostedConfigurationVersionOutput.VersionNumber)
	startDeploymentOutput, err := appConfigAdvance.startDeployment(ctx, configurationProfileId, configurationVersion)
	if err != nil {
		return false, err
	}

	if startDeploymentOutput != nil {
		configurationProfileNameId.Store(configurationName, configurationProfileId)
	}

	return startDeploymentOutput != nil, nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) createConfigurationProfile(ctx context.Context, configurationProfileName string) (*appconfig.CreateConfigurationProfileOutput, error) {
	input := appconfig.CreateConfigurationProfileInput{
		ApplicationId: aws.String(appConfigAdvance.applicationId),
		// 目前只这种类型
		LocationUri: aws.String("hosted"),
		Name:        aws.String(configurationProfileName),
	}
	return appConfigAdvance.appConfigClient.CreateConfigurationProfileWithContext(ctx, &input)
}

func (appConfigAdvance *EnhancedAppConfigAdvance) createHostedConfigurationVersion(ctx context.Context, configurationProfileId string, content string) (*appconfig.CreateHostedConfigurationVersionOutput, error) {
	contentType := http.DetectContentType([]byte(content))
	input := appconfig.CreateHostedConfigurationVersionInput{
		ApplicationId:          aws.String(appConfigAdvance.applicationId),
		ConfigurationProfileId: aws.String(configurationProfileId),
		Content:                []byte(content),
		// text/plain; charset=UTF-8 只取 text/plain
		ContentType: aws.String(strings.SplitN(contentType, "; ", 2)[0]),
	}
	return appConfigAdvance.appConfigClient.CreateHostedConfigurationVersionWithContext(ctx, &input)
}

func (appConfigAdvance *EnhancedAppConfigAdvance) startDeployment(ctx context.Context, configurationProfileId string, configurationVersion string) (*appconfig.StartDeploymentOutput, error) {
	deploymentStrategyId, found := deploymentStrategyNameId.Load(deploymentStrategyName)
	if !found {
		msg := fmt.Sprintf("deploymentStrategy [%s] do not exist in [%s] application", deploymentStrategyName, appConfigAdvance.applicationName)
		return nil, errors.New(msg)
	}
	input := appconfig.StartDeploymentInput{
		ApplicationId:          aws.String(appConfigAdvance.applicationId),
		EnvironmentId:          aws.String(appConfigAdvance.environmentId),
		ConfigurationProfileId: aws.String(configurationProfileId),
		ConfigurationVersion:   aws.String(configurationVersion),
		DeploymentStrategyId:   aws.String(deploymentStrategyId.(string)),
	}
	return appConfigAdvance.appConfigClient.StartDeploymentWithContext(ctx, &input)
}

func (appConfigAdvance *EnhancedAppConfigAdvance) DeleteConfiguration(ctx context.Context, configurationName string) (bool, error) {
	configurationProfileId, found, err := appConfigAdvance.getConfigurationProfileId(ctx, configurationName)
	if err != nil {
		return false, err
	}
	if !found {
		msg := fmt.Sprintf("configuration [%s] do not exist in [%s] environment of [%s] application", configurationName, appConfigAdvance.environmentName, appConfigAdvance.applicationName)
		return false, errors.New(msg)
	}
	output, err := appConfigAdvance.deleteConfigurationProfile(ctx, configurationProfileId)
	if err != nil {
		return false, err
	}

	if output != nil {
		configurationProfileNameId.Delete(configurationName)
	}

	return output != nil, nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) deleteConfigurationProfile(ctx context.Context, configurationProfileId string) (*appconfig.DeleteConfigurationProfileOutput, error) {
	err := appConfigAdvance.deleteAllConfigurationVersion(ctx, configurationProfileId)
	if err != nil {
		return nil, err
	}

	input := appconfig.DeleteConfigurationProfileInput{
		ApplicationId:          aws.String(appConfigAdvance.applicationId),
		ConfigurationProfileId: aws.String(configurationProfileId),
	}
	output, err := appConfigAdvance.appConfigClient.DeleteConfigurationProfileWithContext(ctx, &input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) deleteAllConfigurationVersion(ctx context.Context, configurationProfileId string) error {
	var nextToken *string = nil
	for {
		listHostedConfigurationVersionsInput := appconfig.ListHostedConfigurationVersionsInput{
			ApplicationId:          aws.String(appConfigAdvance.applicationId),
			ConfigurationProfileId: aws.String(configurationProfileId),
			NextToken:              nextToken,
		}
		listHostedConfigurationVersionsOutput, err := appConfigAdvance.appConfigClient.ListHostedConfigurationVersionsWithContext(ctx, &listHostedConfigurationVersionsInput)
		if err != nil {
			return err
		}
		if len(listHostedConfigurationVersionsOutput.Items) == 0 {
			break
		}

		for _, item := range listHostedConfigurationVersionsOutput.Items {
			input := appconfig.DeleteHostedConfigurationVersionInput{
				ApplicationId:          aws.String(appConfigAdvance.applicationId),
				ConfigurationProfileId: item.ConfigurationProfileId,
				VersionNumber:          item.VersionNumber,
			}
			output, err := appConfigAdvance.appConfigClient.DeleteHostedConfigurationVersionWithContext(ctx, &input)
			if err != nil {
				return err
			}
			if output == nil {
				msg := fmt.Sprintf("delete hosted configuration version failed, configurationProfileId: %s, versionNumber: %d", *item.ConfigurationProfileId, *item.VersionNumber)
				return errors.New(msg)
			}
		}

		nextToken = listHostedConfigurationVersionsOutput.NextToken
		if nextToken == nil {
			break
		}
	}
	return nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) nameToId(ctx context.Context) error {
	applicationId, found := applicationNameId[appConfigAdvance.applicationName]
	if !found {
		return fmt.Errorf("can not find application [%s]", appConfigAdvance.applicationName)
	}

	err := appConfigAdvance.listEnvironments(ctx, appConfigAdvance.appConfigClient, applicationId)
	if err != nil {
		return err
	}

	environmentId, found := environmentNameId[appConfigAdvance.environmentName]
	if !found {
		return fmt.Errorf("can not find environment %s at application %s", environmentId, appConfigAdvance.applicationName)
	}

	appConfigAdvance.applicationId = applicationId
	appConfigAdvance.environmentId = environmentId
	return nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) listConfigurationProfiles(ctx context.Context) error {
	var nextToken *string = nil
	for {
		input := appconfig.ListConfigurationProfilesInput{
			ApplicationId: aws.String(appConfigAdvance.applicationId),
			NextToken:     nextToken,
		}
		output, err := appConfigAdvance.appConfigClient.ListConfigurationProfilesWithContext(ctx, &input)
		if err != nil {
			return err
		}
		if len(output.Items) == 0 {
			break
		}

		for _, item := range output.Items {
			configurationProfileNameId.Store(*item.Name, *item.Id)
		}

		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) listDeploymentStrategies(ctx context.Context) error {
	var nextToken *string = nil
	for {
		input := appconfig.ListDeploymentStrategiesInput{
			NextToken: nextToken,
		}
		output, err := appConfigAdvance.appConfigClient.ListDeploymentStrategiesWithContext(ctx, &input)
		if err != nil {
			return err
		}
		if len(output.Items) == 0 {
			break
		}

		for _, item := range output.Items {
			deploymentStrategyNameId.Store(*item.Name, *item.Id)
		}

		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) listApplications(ctx context.Context) error {
	var nextToken *string = nil
	for {
		input := appconfig.ListApplicationsInput{
			NextToken: nextToken,
		}
		output, err := appConfigAdvance.appConfigClient.ListApplicationsWithContext(ctx, &input)
		if err != nil {
			return err
		}
		if len(output.Items) == 0 {
			break
		}

		for _, item := range output.Items {
			applicationNameId[*item.Name] = *item.Id
		}

		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) listEnvironments(ctx context.Context, appConfigClient *appconfig.AppConfig, applicationId string) error {
	var nextToken *string = nil
	for {
		input := appconfig.ListEnvironmentsInput{
			ApplicationId: aws.String(applicationId),
			NextToken:     nextToken,
		}
		output, err := appConfigClient.ListEnvironmentsWithContext(ctx, &input)
		if err != nil {
			return err
		}
		if len(output.Items) == 0 {
			break
		}

		for _, item := range output.Items {
			environmentNameId[*item.Name] = *item.Id
		}

		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return nil
}

func (appConfigAdvance *EnhancedAppConfigAdvance) ApplyOptions(opts ...Option) error {
	for _, opt := range opts {
		err := opt.apply(appConfigAdvance)
		if err != nil {
			return err
		}
	}
	return nil
}
