package constants

import (
	"seall/internal/util"
	"time"
)

const (
	AppName              = "Seall"
	Version              = "3.8.2"
	VersionName          = "Seall"
	GcTime               = time.Minute * 30
	ConfigFileName       = "config.toml"
	MalClientId          = "51cb4294feb400f3ddc66a30f9b9a00f"
	DiscordApplicationId = "1224777421941899285"
	MediaApiUrl          = "https://graphql.simkl.co"
	SimklApiUrl          = "https://api.simkl.com"
	SimklAuthUrl         = "https://simkl.com"
	ProjectRepositoryUrl = "https://github.com/Kolex06/Seall---based-on-seanime"
	ProjectRawMainUrl    = "https://raw.githubusercontent.com/Kolex06/Seall---based-on-seanime/main"
)

const (
	SeallRoomsApiUrl   = ""
	SeallRoomsApiWsUrl = ""
	SeallRoomsVersion  = "1.0.0"
)

var DefaultExtensionMarketplaceURL = ProjectRawMainUrl + "/marketplace.json"
var AnnouncementURL = ProjectRawMainUrl + "/announcements.json"
var InternalMetadataURL = util.Decode("aHR0cHM6Ly9hbmltZS5jbGFwLmluZw==")
