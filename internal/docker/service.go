package docker

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Service struct {
	client *client.Client
}

func NewService() (*Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Service{
		client: cli,
	}, nil
}

func (s *Service) Close() error {
	return s.client.Close()
}

func (s *Service) IsDockerRunning() error {
	_, err := s.client.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return fmt.Errorf("docker is not running or accessible: %w", err)
	}
	return nil
}

func (s *Service) PullImage(ctx context.Context, imageName string) error {
	out, err := s.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageName, err)
	}
	defer out.Close()

	_, err = io.Copy(io.Discard, out)
	if err != nil {
		return fmt.Errorf("failed to read pull response: %w", err)
	}

	return nil
}

type ContainerConfig struct {
	Name          string
	Image         string
	Env           []string
	Ports         map[string]string
	Volumes       []string
	RestartPolicy string
	Public        bool
}

func (s *Service) CreateContainer(ctx context.Context, config *ContainerConfig) (string, error) {
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}

	for containerPort, hostPort := range config.Ports {
		port, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			return "", fmt.Errorf("invalid port %s: %w", containerPort, err)
		}

		hostIP := "127.0.0.1"
		if config.Public {
			hostIP = "0.0.0.0"
		}

		portBindings[port] = []nat.PortBinding{
			{
				HostIP:   hostIP,
				HostPort: hostPort,
			},
		}
		exposedPorts[port] = struct{}{}
	}

	mounts := []mount.Mount{}
	for _, volume := range config.Volumes {
		parts := strings.Split(volume, ":")
		if len(parts) >= 2 {
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: parts[0],
				Target: parts[1],
			})
		}
	}

	restartPolicy := container.RestartPolicyUnlessStopped
	if config.RestartPolicy != "" {
		switch config.RestartPolicy {
		case "no":
			restartPolicy = container.RestartPolicyDisabled
		case "always":
			restartPolicy = container.RestartPolicyAlways
		case "on-failure":
			restartPolicy = container.RestartPolicyOnFailure
		case "unless-stopped":
			restartPolicy = container.RestartPolicyUnlessStopped
		}
	}

	containerConfig := &container.Config{
		Image:        config.Image,
		Env:          config.Env,
		ExposedPorts: exposedPorts,
		Labels: map[string]string{
			"spindb": "true",
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings:  portBindings,
		Mounts:        mounts,
		RestartPolicy: container.RestartPolicy{Name: restartPolicy},
	}

	networkConfig := &network.NetworkingConfig{}

	resp, err := s.client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, config.Name)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return resp.ID, nil
}

func (s *Service) StartContainer(ctx context.Context, containerID string) error {
	return s.client.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (s *Service) StopContainer(ctx context.Context, containerID string) error {
	timeout := 10
	return s.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

func (s *Service) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	return s.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: force})
}

func (s *Service) GetContainer(ctx context.Context, nameOrID string) (*types.ContainerJSON, error) {
	containerList, err := s.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	for _, c := range containerList {
		if c.ID == nameOrID || strings.HasPrefix(c.ID, nameOrID) {
			containerJSON, err := s.client.ContainerInspect(ctx, c.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to inspect container: %w", err)
			}
			return &containerJSON, nil
		}
		for _, name := range c.Names {
			if strings.TrimPrefix(name, "/") == nameOrID {
				containerJSON, err := s.client.ContainerInspect(ctx, c.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to inspect container: %w", err)
				}
				return &containerJSON, nil
			}
		}
	}

	return nil, fmt.Errorf("container %s not found", nameOrID)
}

func (s *Service) IsContainerRunning(ctx context.Context, nameOrID string) (bool, error) {
	containerJSON, err := s.GetContainer(ctx, nameOrID)
	if err != nil {
		return false, err
	}
	return containerJSON.State.Running, nil
}

func (s *Service) WaitForContainer(ctx context.Context, containerID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		running, err := s.IsContainerRunning(ctx, containerID)
		if err != nil {
			return err
		}
		if running {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("container did not start within %v", timeout)
}

func (s *Service) FindAvailablePort(basePort int) (int, error) {
	for port := basePort; port < basePort+100; port++ {
		if s.isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found starting from %d", basePort)
}

func (s *Service) isPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

func (s *Service) GetContainerLogs(ctx context.Context, containerID string, tail string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}

	out, err := s.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer out.Close()

	logs, err := io.ReadAll(out)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(logs), nil
}

func (s *Service) ListSpinDBContainers(ctx context.Context) ([]types.ContainerJSON, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "spindb=true")

	containerList, err := s.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list SpinDB containers: %w", err)
	}

	var containers []types.ContainerJSON
	for _, c := range containerList {
		containerJSON, err := s.client.ContainerInspect(ctx, c.ID)
		if err != nil {
			continue
		}
		containers = append(containers, containerJSON)
	}

	return containers, nil
}

func CreateVolumeMount(hostPath, containerPath string) string {
	return hostPath + ":" + containerPath
}
