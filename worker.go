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

func dctc(client *docker.Client, id string) error {
	if client == nil {
		return errors.New("docker client is nil")
	}

	if id == "" {
		return errors.New("container id is empty")
	}

	container, err := client.InspectContainer(id)
	if err != nil {
		return err
	}

	service := &Service{
		Image:         container.Image,
		ContainerName: strings.TrimLeft(container.Name, "/"),
	}
	if container.HostConfig != nil {
		service.Volumes = container.HostConfig.Binds
		service.Restart = container.HostConfig.RestartPolicy.Name
		service.CapAdd = container.HostConfig.CapAdd
		service.CapDrop = container.HostConfig.CapDrop
		service.CgroupParent = container.HostConfig.CgroupParent
		if len(container.HostConfig.Devices) > 0 {
			service.Devices = make([]string, len(container.HostConfig.Devices))
			for i, device := range container.HostConfig.Devices {
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
		service.DNS = container.HostConfig.DNS
		service.DNSSearch = container.HostConfig.DNSSearch
		service.ExtraHosts = container.HostConfig.ExtraHosts
		service.IPC = container.HostConfig.IpcMode
		service.NetworkMode = container.HostConfig.NetworkMode
		service.Privileged = container.HostConfig.Privileged
		service.Restart = container.HostConfig.RestartPolicy.Name
		service.Sysctls = container.HostConfig.Sysctls
		service.Tmpfs = container.HostConfig.Tmpfs

		if len(container.HostConfig.PortBindings) > 0 {
			service.Ports = make([]string, 0, len(container.HostConfig.PortBindings))
			for port, bindings := range container.HostConfig.PortBindings {
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

		if len(container.HostConfig.Ulimits) > 0 {
			service.Ulimits = make(map[string]*ULimit)
			for _, ulimit := range container.HostConfig.Ulimits {
				service.Ulimits[ulimit.Name] = &ULimit{
					Soft: ulimit.Soft,
					Hard: ulimit.Hard,
				}
			}
		}

	}
	if container.Config != nil {
		service.Environment = container.Config.Env
		if container.Config.Image != "" {
			service.Image = container.Config.Image
		}
		service.Command = container.Config.Cmd
		service.Domainname = container.Config.Domainname
		service.Entrypoint = container.Config.Entrypoint
		if container.Config.Healthcheck != nil {
			service.Healthcheck.Test = container.Config.Healthcheck.Test
			service.Healthcheck.Interval = container.Config.Healthcheck.Interval
			service.Healthcheck.Timeout = container.Config.Healthcheck.Timeout
			service.Healthcheck.Retries = container.Config.Healthcheck.Retries
			service.Healthcheck.StartPeriod = container.Config.Healthcheck.StartPeriod
		}
		service.Labels = container.Config.Labels
		service.MacAddress = container.Config.MacAddress
		service.StdinOpen = container.Config.OpenStdin
		service.StopSignal = container.Config.StopSignal
		service.Tty = container.Config.Tty
		service.User = container.Config.User
		service.WorkingDir = container.Config.WorkingDir

		if len(container.Config.ExposedPorts) > 0 {
			service.Expose = make([]string, 0, len(container.Config.ExposedPorts))
			for port := range container.Config.ExposedPorts {
				service.Expose = append(service.Expose, string(port))
			}
		}

	}
	service.Volumes = make([]string, 0, len(container.Mounts))
	for _, mount := range container.Mounts {
		str := fmt.Sprintf("%s:%s", mount.Source, mount.Destination)
		if mount.Mode != "" {
			str += ":" + mount.Mode
		}
		service.Volumes = append(service.Volumes, str)
	}

	compose := &Compose{
		Version:  "3",
		Services: make(map[string]*Service),
	}
	compose.Services[service.ContainerName] = service

	data, err := yml.Marshal(compose)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
