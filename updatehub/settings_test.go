/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier:     GPL-2.0
 */

package updatehub

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const customSettings = `
[Polling]
Interval=1
Enabled=false
LastPoll=2017-01-01T00:00:00Z
FirstPoll=2017-02-02T00:00:00Z
ExtraInterval=4
Retries=5

[Storage]
ReadOnly=true

[Update]
DownloadDir=/tmp/download
AutoDownloadWhenAvailable=false
AutoInstallAfterDownload=false
AutoRebootAfterInstall=false
SupportedInstallModes=mode1,mode2

[Network]
DisableHttps=true
UpdateHubServerAddress=localhost

[Firmware]
MetadataPath=/tmp/metadata
`

func TestToString(t *testing.T) {
	s := &Settings{
		PollingSettings: PollingSettings{
			PollingInterval: defaultPollingInterval,
			PollingEnabled:  true,
			PersistentPollingSettings: PersistentPollingSettings{
				LastPoll:             (time.Time{}).UTC(),
				FirstPoll:            (time.Time{}).UTC(),
				ExtraPollingInterval: 0,
				PollingRetries:       0,
			},
		},

		StorageSettings: StorageSettings{
			ReadOnly: false,
		},

		UpdateSettings: UpdateSettings{
			DownloadDir:               "/tmp",
			AutoDownloadWhenAvailable: true,
			AutoInstallAfterDownload:  true,
			AutoRebootAfterInstall:    true,
			SupportedInstallModes:     []string{"dry-run", "copy", "flash", "imxkobs", "raw", "tarball", "ubifs"},
		},

		NetworkSettings: NetworkSettings{
			DisableHTTPS:  false,
			ServerAddress: "",
		},

		FirmwareSettings: FirmwareSettings{
			FirmwareMetadataPath: "",
		},
	}

	outputJSON, _ := json.MarshalIndent(s, "", "    ")
	expectedString := string(outputJSON)

	assert.Equal(t, expectedString, s.ToString())
}

func TestLoadSettingsDefaultValues(t *testing.T) {
	s, err := LoadSettings(bytes.NewReader([]byte("")))
	assert.NoError(t, err)

	assert.Equal(t, time.Hour, s.PollingSettings.PollingInterval)
	assert.Equal(t, true, s.PollingSettings.PollingEnabled)
	assert.Equal(t, (time.Time{}).UTC(), s.PollingSettings.PersistentPollingSettings.LastPoll)
	assert.Equal(t, (time.Time{}).UTC(), s.PollingSettings.PersistentPollingSettings.FirstPoll)
	assert.Equal(t, time.Duration(0), s.PollingSettings.PersistentPollingSettings.ExtraPollingInterval)
	assert.Equal(t, 0, s.PollingSettings.PersistentPollingSettings.PollingRetries)

	assert.Equal(t, false, s.StorageSettings.ReadOnly)

	assert.Equal(t, "/tmp", s.UpdateSettings.DownloadDir)
	assert.Equal(t, true, s.UpdateSettings.AutoDownloadWhenAvailable)
	assert.Equal(t, true, s.UpdateSettings.AutoInstallAfterDownload)
	assert.Equal(t, true, s.UpdateSettings.AutoRebootAfterInstall)
	assert.Equal(t, []string{"dry-run", "copy", "flash", "imxkobs", "raw", "tarball", "ubifs"}, s.UpdateSettings.SupportedInstallModes)

	assert.Equal(t, false, s.NetworkSettings.DisableHTTPS)
	assert.Equal(t, "", s.NetworkSettings.ServerAddress)

	assert.Equal(t, "", s.FirmwareSettings.FirmwareMetadataPath)
}

func TestLoadSettings(t *testing.T) {
	testCases := []struct {
		name             string
		data             string
		expectedSettings *Settings
	}{
		{
			"DefaultValues",
			"",
			&Settings{
				PollingSettings: PollingSettings{
					PollingInterval: defaultPollingInterval,
					PollingEnabled:  true,
					PersistentPollingSettings: PersistentPollingSettings{
						LastPoll:             (time.Time{}).UTC(),
						FirstPoll:            (time.Time{}).UTC(),
						ExtraPollingInterval: 0,
						PollingRetries:       0,
					},
				},

				StorageSettings: StorageSettings{
					ReadOnly: false,
				},

				UpdateSettings: UpdateSettings{
					DownloadDir:               "/tmp",
					AutoDownloadWhenAvailable: true,
					AutoInstallAfterDownload:  true,
					AutoRebootAfterInstall:    true,
					SupportedInstallModes:     []string{"dry-run", "copy", "flash", "imxkobs", "raw", "tarball", "ubifs"},
				},

				NetworkSettings: NetworkSettings{
					DisableHTTPS:  false,
					ServerAddress: "",
				},

				FirmwareSettings: FirmwareSettings{
					FirmwareMetadataPath: "",
				},
			},
		},

		{
			"CustomValues",
			customSettings,
			&Settings{
				PollingSettings: PollingSettings{
					PollingInterval: 1,
					PollingEnabled:  false,
					PersistentPollingSettings: PersistentPollingSettings{
						LastPoll:             time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
						FirstPoll:            time.Date(2017, time.February, 2, 0, 0, 0, 0, time.UTC),
						ExtraPollingInterval: 4,
						PollingRetries:       5,
					},
				},

				StorageSettings: StorageSettings{
					ReadOnly: true,
				},

				UpdateSettings: UpdateSettings{
					DownloadDir:               "/tmp/download",
					AutoDownloadWhenAvailable: false,
					AutoInstallAfterDownload:  false,
					AutoRebootAfterInstall:    false,
					SupportedInstallModes:     []string{"mode1", "mode2"},
				},

				NetworkSettings: NetworkSettings{
					DisableHTTPS:  true,
					ServerAddress: "localhost",
				},

				FirmwareSettings: FirmwareSettings{
					FirmwareMetadataPath: "/tmp/metadata",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := LoadSettings(bytes.NewReader([]byte(tc.data)))
			assert.NoError(t, err)
			assert.NotNil(t, s)
			assert.Equal(t, tc.expectedSettings, s)
		})
	}
}

func TestSaveSettings(t *testing.T) {
	expectedSettings, err := LoadSettings(bytes.NewReader([]byte("")))
	assert.NoError(t, err)
	assert.NotNil(t, expectedSettings)

	var buf bytes.Buffer
	err = SaveSettings(expectedSettings, &buf)
	assert.NoError(t, err)

	s, err := LoadSettings(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.Equal(t, expectedSettings, s)
}