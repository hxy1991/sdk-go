package awsappconfigdata

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnhancedAppConfig_getConfiguration(t *testing.T) {
	appConfig, err := NewWithOptions(
		context.TODO(),
		WithCacheRefreshInterval(time.Second*10),
		WithApplicationName("APP"),
	)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	got, err := appConfig.getConfiguration(context.TODO(), "version.json", nil)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	t.Log(string(got.Configuration))
	t.Log(got)

	time.Sleep(time.Second * 15)

	got1, err := appConfig.getConfiguration(context.TODO(), "version.json", got.NextPollConfigurationToken)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	assert.Equal(t, "", string(got1.Configuration))
}
