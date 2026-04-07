package convert

import (
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func FromExecdFileInfo(info execdapi.FileInfo) models.EntryInfo {
	return models.EntryInfo{
		Path:       info.Path,
		Size:       info.Size,
		ModifiedAt: info.ModifiedAt,
		CreatedAt:  info.CreatedAt,
		Mode:       info.Mode,
		Owner:      info.Owner,
		Group:      info.Group,
	}
}

func FromExecdFileInfoMap(items *map[string]execdapi.FileInfo) map[string]models.EntryInfo {
	if items == nil {
		return map[string]models.EntryInfo{}
	}

	out := make(map[string]models.EntryInfo, len(*items))
	for path, info := range *items {
		out[path] = FromExecdFileInfo(info)
	}
	return out
}

func FromExecdFileInfoList(items *[]execdapi.FileInfo) []models.EntryInfo {
	if items == nil {
		return []models.EntryInfo{}
	}

	out := make([]models.EntryInfo, 0, len(*items))
	for _, item := range *items {
		out = append(out, FromExecdFileInfo(item))
	}
	return out
}

func ToExecdPermission(mode int, owner, group string) execdapi.Permission {
	permission := execdapi.Permission{Mode: mode}
	if owner != "" {
		permission.Owner = &owner
	}
	if group != "" {
		permission.Group = &group
	}
	return permission
}
