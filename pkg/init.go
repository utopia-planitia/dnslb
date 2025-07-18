package dnslb

import (
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go/v4"
)

func initAPI(apiToken, apiKey, apiEmail string) (api *cloudflare.API, err error) {
	if apiToken == "" && (apiKey == "" || apiEmail == "") {
		return nil, fmt.Errorf("either CF_API_TOKEN or CF_API_KEY and CF_API_EMAIL need to be set")
	}

	if apiToken != "" {
		return cloudflare.NewWithAPIToken(apiToken)
	}

	return cloudflare.New(apiKey, apiEmail)
}

func loadZoneID(api *cloudflare.API) (string, string, error) {
	zone := os.Getenv("CF_ZONE")
	if zone == "" {
		return "", "", fmt.Errorf("CF_ZONE needs to be set")
	}

	zoneID, err := api.ZoneIDByName(zone)
	if err != nil {
		return "", "", err
	}

	return zone, zoneID, nil
}
