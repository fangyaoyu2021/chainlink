package periodicbackup

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriodicBackup_RunBackup(t *testing.T) {
	rawConfig := orm.NewConfig()
	backupConfig := newTestConfig(time.Minute, nil, rawConfig.DatabaseURL(), os.TempDir())
	periodicBackup := NewDatabaseBackup(backupConfig, logger.Default).(*databaseBackup)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup("0.9.9")
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Greater(t, file.Size(), int64(0))
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "cl_backup_0.9.9")
}

func TestPeriodicBackup_RunBackupWithoutVersion(t *testing.T) {
	rawConfig := orm.NewConfig()
	backupConfig := newTestConfig(time.Minute, nil, rawConfig.DatabaseURL(), os.TempDir())
	periodicBackup := NewDatabaseBackup(backupConfig, logger.Default).(*databaseBackup)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup("unset")
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Greater(t, file.Size(), int64(0))
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "cl_backup_unset")
}

func TestPeriodicBackup_RunBackupViaAltUrl(t *testing.T) {
	rawConfig := orm.NewConfig()
	altUrl, _ := url.Parse("postgresql//invalid")
	backupConfig := newTestConfig(time.Minute, altUrl, rawConfig.DatabaseURL(), os.TempDir())
	periodicBackup := NewDatabaseBackup(backupConfig, logger.Default).(*databaseBackup)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	_, err := periodicBackup.runBackup("")
	require.Error(t, err, "connection to database \"postgresql//invalid\" failed")
}

func TestPeriodicBackup_FrequencyTooSmall(t *testing.T) {
	rawConfig := orm.NewConfig()
	backupConfig := newTestConfig(time.Second, nil, rawConfig.DatabaseURL(), os.TempDir())
	periodicBackup := NewDatabaseBackup(backupConfig, logger.Default).(*databaseBackup)
	assert.True(t, periodicBackup.frequencyIsTooSmall())
}

type testConfig struct {
	databaseBackupFrequency time.Duration
	databaseBackupURL       *url.URL
	databaseURL             url.URL
	rootDir                 string
}

func (config testConfig) DatabaseBackupFrequency() time.Duration {
	return config.databaseBackupFrequency
}
func (config testConfig) DatabaseBackupURL() *url.URL {
	return config.databaseBackupURL
}
func (config testConfig) DatabaseURL() url.URL {
	return config.databaseURL
}
func (config testConfig) RootDir() string {
	return config.rootDir
}

func newTestConfig(frequency time.Duration, databaseBackupURL *url.URL, databaseURL url.URL, outputParentDir string) testConfig {
	return testConfig{
		databaseBackupFrequency: frequency,
		databaseBackupURL:       databaseBackupURL,
		databaseURL:             databaseURL,
		rootDir:                 outputParentDir,
	}
}
