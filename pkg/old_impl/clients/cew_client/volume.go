package cew_client

import (
	"context"
	"mgw-module-manager-migration/pkg/old_impl/libs/cew_lib"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetVolumes(ctx context.Context, filter cew_lib.VolumeFilter) ([]cew_lib.Volume, error) {
	u, err := url.JoinPath(c.baseUrl, cew_lib.VolumesPath)
	if err != nil {
		return nil, err
	}
	u += genVolumesQuery(filter)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var volumes []cew_lib.Volume
	err = c.baseClient.ExecRequestJSON(req, &volumes)
	if err != nil {
		return nil, err
	}
	return volumes, nil
}

func (c *Client) GetVolume(ctx context.Context, id string) (cew_lib.Volume, error) {
	u, err := url.JoinPath(c.baseUrl, cew_lib.VolumesPath, id)
	if err != nil {
		return cew_lib.Volume{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return cew_lib.Volume{}, err
	}
	var volume cew_lib.Volume
	err = c.baseClient.ExecRequestJSON(req, &volume)
	if err != nil {
		return cew_lib.Volume{}, err
	}
	return volume, nil
}

func genVolumesQuery(filter cew_lib.VolumeFilter) string {
	var q []string
	if len(filter.Labels) > 0 {
		q = append(q, "labels="+genLabels(filter.Labels, "=", ","))
	}
	if len(q) > 0 {
		return "?" + strings.Join(q, "&")
	}
	return ""
}
