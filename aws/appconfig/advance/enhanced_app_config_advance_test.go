package awsappconfigadvance

import (
	"context"
	"fmt"
	"github.com/hxy1991/sdk-go/constant"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hxy1991/sdk-go/aws/appconfig"
	"github.com/stretchr/testify/assert"
)

var regionName = "us-east-1"
var applicationName = "app1"
var environmentName = "Test"

func TestAppConfigAdvance(t *testing.T) {
	setEnvs(t)

	ctx := context.Background()

	appConfig, err := awsappconfig.NewWithOptions(
		awsappconfig.WithApplicationName(applicationName),
		awsappconfig.WithCacheRefreshInterval(time.Second*10),
	)
	assert.Nil(t, err)

	appConfigAdvance, err := NewWithApplicationName(applicationName)
	assert.Nil(t, err)

	configurationName := fmt.Sprintf("TestAppConfigAdvance-%d", time.Now().Unix())
	testCreate(ctx, t, appConfig, appConfigAdvance, configurationName)
	testUpdate(ctx, t, appConfig, appConfigAdvance, configurationName)
	testDelete(ctx, t, appConfig, appConfigAdvance, configurationName)
}

func testDelete(ctx context.Context, t *testing.T, appConfig *awsappconfig.EnhancedAppConfig, appConfigAdvance *EnhancedAppConfigAdvance, configurationName string) {
	t.Log("start testDelete")
	deleteAppConfig(ctx, t, appConfigAdvance, configurationName)
	now := time.Now()
	for i := 1; ; i++ {
		_, err := appConfig.GetConfiguration(ctx, configurationName)
		if err == nil {
			//t.Log("expected err but got nil")
		} else {
			msg := fmt.Sprintf("Configuration Profile Id %s could not be found", configurationName)
			if strings.Contains(err.Error(), msg) {
				t2 := time.Now()
				// 删除需要4秒多才能获取
				t.Log("count: ", i, ", cost: ", t2.Sub(now))
				break
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
	t.Log("end testDelete")
}

func testUpdate(ctx context.Context, t *testing.T, appConfig *awsappconfig.EnhancedAppConfig, appConfigAdvance *EnhancedAppConfigAdvance, configurationName string) time.Time {
	t.Log("start testUpdate")
	updateContent := update(ctx, t, appConfigAdvance, configurationName)
	now := time.Now()
	for i := 1; ; i++ {
		getContent := get(ctx, t, appConfig, configurationName)
		if getContent != updateContent {
			//t.Log("expected ", updateContent, " but got ", getContent)
		} else {
			t2 := time.Now()
			// 更新需要4秒多才能获取到
			t.Log("count: ", i, ", cost: ", t2.Sub(now), " updateContent == getContent == ", getContent)
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
	t.Log("end testUpdate")
	return now
}

func testCreate(ctx context.Context, t *testing.T, appConfig *awsappconfig.EnhancedAppConfig, appConfigAdvance *EnhancedAppConfigAdvance, configurationName string) time.Time {
	t.Log("start testCreate")
	createContent := create(ctx, t, appConfigAdvance, configurationName)
	now := time.Now()
	for i := 1; ; i++ {
		getContent := get(ctx, t, appConfig, configurationName)
		if createContent != getContent {
			t.Log("expected ", createContent, " but got ", getContent)
		} else {
			t.Log("count: ", i, ", cost: ", time.Since(now), " createContent == getContent == ", getContent)
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
	t.Log("end testCreate")
	return now
}

func create(ctx context.Context, t *testing.T, appConfigAdvance *EnhancedAppConfigAdvance, configurationName string) string {
	content := time.Now().Format(time.RFC3339)
	isSuccess, err := appConfigAdvance.CreateConfiguration(ctx, configurationName, content)
	assert.Nil(t, err)
	assert.True(t, isSuccess)
	return content

}
func get(ctx context.Context, t *testing.T, appConfig *awsappconfig.EnhancedAppConfig, configurationName string) string {
	content, err := appConfig.GetConfiguration(ctx, configurationName)
	assert.Nil(t, err)
	return content
}

func update(ctx context.Context, t *testing.T, appConfigAdvance *EnhancedAppConfigAdvance, configurationName string) string {
	updateContent := time.Now().Format(time.RFC3339)
	isSuccess, err := appConfigAdvance.UpdateConfiguration(ctx, configurationName, updateContent)
	assert.Nil(t, err)
	assert.True(t, isSuccess)
	return updateContent
}

func deleteAppConfig(ctx context.Context, t *testing.T, appConfigAdvance *EnhancedAppConfigAdvance, configurationName string) {
	isSuccess, err := appConfigAdvance.DeleteConfiguration(ctx, configurationName)
	assert.Nil(t, err)
	assert.True(t, isSuccess)
}

func setEnvs(t *testing.T) {
	setEnv(t, constant.RegionEnvName, regionName)
	setEnv(t, constant.EnvironmentEnvName, environmentName)
}

func setEnv(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	assert.Nil(t, err)
}
