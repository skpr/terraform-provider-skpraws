package types

import (
	awstypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func AssetCategoryTypeFromString(s string) awstypes.AssetCategoryType {
	mapping := map[string]awstypes.AssetCategoryType{
		"PAGE_BACKGROUND": awstypes.AssetCategoryTypePageBackground,
	}
	return mapping[s]
}

func ColorSchemeModeTypeFromString(s string) awstypes.ColorSchemeModeType {
	mapping := map[string]awstypes.ColorSchemeModeType{
		"LIGHT":   awstypes.ColorSchemeModeTypeLight,
		"DARK":    awstypes.ColorSchemeModeTypeDark,
		"DYNAMIC": awstypes.ColorSchemeModeTypeDynamic,
	}
	return mapping[s]
}

func AssetExtensionTypeFromString(s string) awstypes.AssetExtensionType {
	mapping := map[string]awstypes.AssetExtensionType{
		"ICO":  awstypes.AssetExtensionTypeIco,
		"JPEG": awstypes.AssetExtensionTypeJpeg,
		"PNG":  awstypes.AssetExtensionTypePng,
		"SVG":  awstypes.AssetExtensionTypeSvg,
		"WEBP": awstypes.AssetExtensionTypeWebp,
	}
	return mapping[s]
}
