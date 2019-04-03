package main

import (
	"errors"
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	yml "gopkg.in/yaml.v2"
)

func newClient() (*docker.Client, error) {
	var client *docker.Client
	var err error
	if tls {
		client, err = docker.NewTLSClient(
			host,
			tlscert,
			tlskey,
			tlscacert)
	} else {
		client, err = docker.NewClient(host)
	}
	return client, err
}

func dctc(client *docker.Client, ids []string) (string, error) {
	if client == nil {
		return "", errors.New("docker client is nil")
	}

	if ids == nil || len(ids) == 0 {
		return "", errors.New("container id is empty")
	}

	compose := &Compose{
		Version:  "3",
		Services: make(map[string]*Service),
	}

	exists := make(map[string]bool)

	for _, id := range ids {
		if exists[id] {
			continue
		}

		exists[id] = true

		container, err := client.InspectContainer(id)
		if err != nil {
			return "", err
		}

		service := &Service{
			Image:         container.Image,
			ContainerName: strings.TrimLeft(container.Name, "/"),
		}
		if container.HostConfig != nil {
			convertHostConfig(service, container.HostConfig)
		}
		if container.Config != nil {
			convertConfig(service, container.Config)
		}
		service.Volumes = make([]string, 0, len(container.Mounts))
		for _, mount := range container.Mounts {
			str := fmt.Sprintf("%s:%s", mount.Source, mount.Destination)
			if mount.Mode != "" {
				str += ":" + mount.Mode
			}
			service.Volumes = append(service.Volumes, str)
		}

		compose.Services[service.ContainerName] = service
	}

	data, err := yml.Marshal(compose)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func convertHostConfig(service *Service, hostConfig *docker.HostConfig) {
	if service == nil || hostConfig == nil {
		return
	}

	service.Volumes = hostConfig.Binds
	service.Restart = hostConfig.RestartPolicy.Name
	service.CapAdd = hostConfig.CapAdd
	service.CapDrop = hostConfig.CapDrop
	service.CgroupParent = hostConfig.CgroupParent
	if len(hostConfig.Devices) > 0 {
		service.Devices = make([]string, len(hostConfig.Devices))
		for i, device := range hostConfig.Devices {
			strs := make([]string, 0, 3)
			if device.PathOnHost != "" {
				strs = append(strs, device.PathOnHost)
			}
			if device.PathInContainer != "" {
				strs = append(strs, device.PathInContainer)
			}
			if device.CgroupPermissions != "" {
				strs = append(strs, device.CgroupPermissions)
			}
			service.Devices[i] = strings.Join(strs, ":")
		}
	}
	service.DNS = hostConfig.DNS
	service.DNSSearch = hostConfig.DNSSearch
	service.ExtraHosts = hostConfig.ExtraHosts
	service.IPC = hostConfig.IpcMode
	service.NetworkMode = hostConfig.NetworkMode
	service.Privileged = hostConfig.Privileged
	service.Restart = hostConfig.RestartPolicy.Name
	service.Sysctls = hostConfig.Sysctls
	service.Tmpfs = hostConfig.Tmpfs

	if len(hostConfig.PortBindings) > 0 {
		service.Ports = make([]string, 0, len(hostConfig.PortBindings))
		for port, bindings := range hostConfig.PortBindings {
			for _, binding := range bindings {
				strs := make([]string, 0, 3)
				if binding.HostIP != "" && binding.HostIP != "0.0.0.0" {
					strs = append(strs, binding.HostIP)
				}
				if binding.HostPort != "" {
					strs = append(strs, binding.HostPort)
				}
				strs = append(strs, string(port))
				service.Ports = append(service.Ports, strings.Join(strs, ":"))
			}
		}
	}

	if len(hostConfig.Ulimits) > 0 {
		service.Ulimits = make(map[string]*ULimit)
		for _, ulimit := range hostConfig.Ulimits {
			service.Ulimits[ulimit.Name] = &ULimit{
				Soft: ulimit.Soft,
				Hard: ulimit.Hard,
			}
		}
	}
}

func convertConfig(service *Service, config *docker.Config) {
	if service == nil || config == nil {
		return
	}
	service.Environment = config.Env
	if config.Image != "" {
		service.Image = config.Image
	}
	service.Command = config.Cmd
	service.Domainname = config.Domainname
	service.Entrypoint = config.Entrypoint
	if config.Healthcheck != nil {
		service.Healthcheck.Test = config.Healthcheck.Test
		service.Healthcheck.Interval = config.Healthcheck.Interval
		service.Healthcheck.Timeout = config.Healthcheck.Timeout
		service.Healthcheck.Retries = config.Healthcheck.Retries
		service.Healthcheck.StartPeriod = config.Healthcheck.StartPeriod
	}
	service.Labels = config.Labels
	service.MacAddress = config.MacAddress
	service.StdinOpen = config.OpenStdin
	service.StopSignal = config.StopSignal
	service.Tty = config.Tty
	service.User = config.User
	service.WorkingDir = config.WorkingDir

	if len(config.ExposedPorts) > 0 {
		service.Expose = make([]string, 0, len(config.ExposedPorts))
		for port := range config.ExposedPorts {
			service.Expose = append(service.Expose, string(port))
		}
	}
}
