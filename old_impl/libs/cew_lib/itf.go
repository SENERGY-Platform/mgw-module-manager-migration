package cew_lib

import (
	"context"
	"io"
)

type Api interface {
	GetContainers(ctx context.Context, filter ContainerFilter) ([]Container, error)
	GetContainer(ctx context.Context, id string) (Container, error)
	CreateContainer(ctx context.Context, container Container) (id string, err error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) (jobId string, err error)
	RestartContainer(ctx context.Context, id string) (jobId string, err error)
	RemoveContainer(ctx context.Context, id string, force bool) error
	GetContainerLog(ctx context.Context, id string, logOptions LogFilter) (io.ReadCloser, error)
	ContainerExec(ctx context.Context, id string, exeConf ExecConfig) (string, error)
	GetImages(ctx context.Context, filter ImageFilter) ([]Image, error)
	GetImage(ctx context.Context, id string) (Image, error)
	AddImage(ctx context.Context, img string) (jobId string, err error)
	RemoveImage(ctx context.Context, id string) error
	GetNetworks(ctx context.Context) ([]Network, error)
	GetNetwork(ctx context.Context, id string) (Network, error)
	CreateNetwork(ctx context.Context, net Network) (string, error)
	RemoveNetwork(ctx context.Context, id string) error
	GetVolumes(ctx context.Context, filter VolumeFilter) ([]Volume, error)
	GetVolume(ctx context.Context, id string) (Volume, error)
	CreateVolume(ctx context.Context, vol Volume) (string, error)
	RemoveVolume(ctx context.Context, id string, force bool) error
}
