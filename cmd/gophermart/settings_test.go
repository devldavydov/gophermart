package main

import (
	"flag"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceSettingsAdaptDefault(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	serviceSettings, err := ServiceSettingsAdapt(config)
	assert.NoError(t, err)

	expRunAddress, _ := url.Parse("http://127.0.0.1:8080")
	expAccrAddress, _ := url.Parse("http://127.0.0.1:9090")
	assert.Equal(t, expRunAddress, serviceSettings.RunAddress)
	assert.Equal(t, "", serviceSettings.DatabaseDsn)
	assert.Equal(t, expAccrAddress, serviceSettings.AccrualSystemAddress)
	assert.Equal(t, 10*time.Second, serviceSettings.ShutdownTimeout)
}

func TestServiceSettingsAdaptCustomEnv(t *testing.T) {
	t.Setenv("RUN_ADDRESS", "1.1.1.1:9999")
	t.Setenv("DATABASE_URI", "postgre:1234")
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "2.2.2.2:9999")
	t.Setenv("SHUTDOWN_TIMEOUT", "5s")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	serviceSettings, err := ServiceSettingsAdapt(config)
	assert.NoError(t, err)

	expRunAddress, _ := url.Parse("http://1.1.1.1:9999")
	expAccrAddress, _ := url.Parse("http://2.2.2.2:9999")
	assert.Equal(t, expRunAddress, serviceSettings.RunAddress)
	assert.Equal(t, "postgre:1234", serviceSettings.DatabaseDsn)
	assert.Equal(t, expAccrAddress, serviceSettings.AccrualSystemAddress)
	assert.Equal(t, 5*time.Second, serviceSettings.ShutdownTimeout)
}

func TestServiceSettingsAdaptCustomFlag(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "1.1.1.1:9999", "-d", "postgre:1234", "-r", "2.2.2.2:9999", "-t", "5s"})
	assert.NoError(t, err)

	serviceSettings, err := ServiceSettingsAdapt(config)
	assert.NoError(t, err)

	expRunAddress, _ := url.Parse("http://1.1.1.1:9999")
	expAccrAddress, _ := url.Parse("http://2.2.2.2:9999")
	assert.Equal(t, expRunAddress, serviceSettings.RunAddress)
	assert.Equal(t, "postgre:1234", serviceSettings.DatabaseDsn)
	assert.Equal(t, expAccrAddress, serviceSettings.AccrualSystemAddress)
	assert.Equal(t, 5*time.Second, serviceSettings.ShutdownTimeout)
}

func TestServiceSettingsAdaptCustomEnvAndFlag(t *testing.T) {
	t.Setenv("RUN_ADDRESS", "3.3.3.3:9999")
	t.Setenv("DATABASE_URI", "postgre:4567")
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "3.3.3.3:9999")
	t.Setenv("SHUTDOWN_TIMEOUT", "50s")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "1.1.1.1:9999", "-d", "postgre:1234", "-r", "2.2.2.2:9999", "-t", "5s"})
	assert.NoError(t, err)

	serviceSettings, err := ServiceSettingsAdapt(config)
	assert.NoError(t, err)

	expRunAddress, _ := url.Parse("http://3.3.3.3:9999")
	expAccrAddress, _ := url.Parse("http://3.3.3.3:9999")
	assert.Equal(t, expRunAddress, serviceSettings.RunAddress)
	assert.Equal(t, "postgre:4567", serviceSettings.DatabaseDsn)
	assert.Equal(t, expAccrAddress, serviceSettings.AccrualSystemAddress)
	assert.Equal(t, 50*time.Second, serviceSettings.ShutdownTimeout)
}

func TestServiceSettingsAdaptCustomEnvAndFlagMix(t *testing.T) {
	t.Setenv("RUN_ADDRESS", "3.3.3.3:9999")
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "3.3.3.3:9999")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "1.1.1.1:9999", "-d", "postgre:4567"})
	assert.NoError(t, err)

	serviceSettings, err := ServiceSettingsAdapt(config)
	assert.NoError(t, err)

	expRunAddress, _ := url.Parse("http://3.3.3.3:9999")
	expAccrAddress, _ := url.Parse("http://3.3.3.3:9999")
	assert.Equal(t, expRunAddress, serviceSettings.RunAddress)
	assert.Equal(t, "postgre:4567", serviceSettings.DatabaseDsn)
	assert.Equal(t, expAccrAddress, serviceSettings.AccrualSystemAddress)
	assert.Equal(t, 10*time.Second, serviceSettings.ShutdownTimeout)
}

func TestServiceSettingsAdaptCustomError(t *testing.T) {
	for _, envVar := range []string{"RUN_ADDRESS", "ACCRUAL_SYSTEM_ADDRESS"} {
		t.Run(envVar, func(t *testing.T) {
			t.Setenv(envVar, "a.%^7b.c.d.e.f")

			testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
			config, err := LoadConfig(*testFlagSet, []string{})
			assert.NoError(t, err)

			_, err = ServiceSettingsAdapt(config)
			assert.Error(t, err)
		})
	}
}
