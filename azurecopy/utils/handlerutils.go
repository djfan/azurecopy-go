package utils

import (
	"azurecopy/azurecopy/handlers"
	"azurecopy/azurecopy/models"
	"azurecopy/azurecopy/utils/misc"
	"crypto/md5"
	"encoding/hex"

	log "github.com/Sirupsen/logrus"
)

// GetHandler gets the appropriate handler for the cloudtype.
// Should I be doing this another way?
func GetHandler(cloudType models.CloudType, isSource bool, config misc.CloudConfig, cacheToDisk bool, isEmulator bool) handlers.CloudHandlerInterface {
	switch cloudType {
	case models.Azure:

		accountName, accountKey := getAzureCredentials(isSource, config)

		log.Debug("Got Azure Handler")
		ah, _ := handlers.NewAzureHandler(accountName, accountKey, isSource, cacheToDisk, isEmulator)
		return ah

	case models.Filesystem:
		log.Debug("Got Filesystem Handler")
		fh, _ := handlers.NewFilesystemHandler("c:/temp/", isSource) // default path?
		return fh

	case models.S3:
		log.Debug("Got S3 Handler")
		accessID, accessSecret, region := getS3Credentials(isSource, config)

		sh, _ := handlers.NewS3Handler(accessID, accessSecret, region, isSource, true)
		return sh

	}

	return nil
}

func GenerateCacheName(path string) string {
	hasher := md5.New()
	hasher.Write([]byte(path))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getAzureCredentials(isSource bool, config misc.CloudConfig) (accountName string, accountKey string) {
	if isSource {
		accountName = config.Configuration[misc.AzureSourceAccountName]
		accountKey = config.Configuration[misc.AzureSourceAccountKey]
	} else {
		accountName = config.Configuration[misc.AzureDestAccountName]
		accountKey = config.Configuration[misc.AzureDestAccountKey]
	}

	if accountName == "" || accountKey == "" {
		accountName = config.Configuration[misc.AzureDefaultAccountName]
		accountKey = config.Configuration[misc.AzureDefaultAccountKey]
	}

	return accountName, accountKey
}

func getS3Credentials(isSource bool, config misc.CloudConfig) (accessID string, accessSecret string, region string) {
	if isSource {
		accessID = config.Configuration[misc.S3SourceAccessID]
		accessSecret = config.Configuration[misc.S3SourceAccessSecret]
		region = config.Configuration[misc.S3SourceRegion]
	} else {
		accessID = config.Configuration[misc.S3DestAccessID]
		accessSecret = config.Configuration[misc.S3DestAccessSecret]
		region = config.Configuration[misc.S3DestRegion]
	}

	if accessID == "" || accessSecret == "" {
		accessID = config.Configuration[misc.S3DefaultAccessID]
		accessSecret = config.Configuration[misc.S3DefaultAccessSecret]
		region = config.Configuration[misc.S3DefaultRegion]
	}

	return accessID, accessSecret, region
}
