package main

import "time"

// Service service config
type Service struct {
	CapAdd        []string `yaml:"cap_add,omitempty"`
	CapDrop       []string `yaml:"cap_drop,omitempty"`
	CgroupParent  string   `yaml:"cgroup_parent,omitempty"`
	Command       []string `yaml:"command,omitempty"`
	ContainerName string   `yaml:"container_name,omitempty"`
	Devices       []string `yaml:"devices,omitempty"`
	DNS           []string `yaml:"dns,omitempty"`
	DNSSearch     []string `yaml:"dns_search,omitempty"`
	Domainname    string   `yaml:"domainname,omitempty"`
	Entrypoint    []string `yaml:"entrypoint,omitempty"`
	Environment   []string `yaml:"environment,omitempty"`
	Expose        []string `yaml:"expose,omitempty"`
	ExtraHosts    []string `yaml:"extra_hosts,omitempty"`
	Healthcheck   struct {
		Test        []string      `yaml:"test,omitempty"`
		Interval    time.Duration `yaml:"interval,omitempty"`
		Timeout     time.Duration `yaml:"timeout,omitempty"`
		Retries     int           `yaml:"retries,omitempty"`
		StartPeriod time.Duration `yaml:"start_period,omitempty"`
	} `yaml:"healcheck,omitempty"`
	Image       string            `yaml:"image,omitempty"`
	IPC         string            `yaml:"ipc,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Links       []string          `yaml:"links,omitempty"`
	MacAddress  string            `yaml:"mac_address,omitempty"`
	NetworkMode string            `yaml:"network_mode,omitempty"`
	Networks    map[string]struct {
		Aliases     []string `yaml:"aliases,omitempty"`
		IPv4Address string   `yaml:"ipv4_address,omitempty"`
		IPv6Address string   `yaml:"ipv6_address,omitempty"`
	} `yaml:"networks,omitempty"` // TODO
	Pid             string             `yaml:"pid,omitempty"` // TODO
	Ports           []string           `yaml:"ports,omitempty"`
	Privileged      bool               `yaml:"privileged,omitempty"`
	ReadOnly        bool               `yaml:"read_only,omitempty"` // TODO
	Restart         string             `yaml:"restart,omitempty"`
	StdinOpen       bool               `yaml:"stdin_open,omitempty"`
	StopGracePeriod string             `yaml:"stop_grace_period,omitempty"`
	StopSignal      string             `yaml:"stop_signal,omitempty"`
	Sysctls         map[string]string  `yaml:"sysctls,omitempty"`
	Tmpfs           map[string]string  `yaml:"tmpfs,omitempty"`
	Tty             bool               `yaml:"tty,omitempty"`
	Ulimits         map[string]*ULimit `yaml:"ulimits,omitempty"`
	User            string             `yaml:"user,omitempty"`
	Volumes         []string           `yaml:"volumes,omitempty"` // TODO
	WorkingDir      string             `yaml:"working_dir,omitempty"`
}

// ULimit defines system-wide resource limitations This can help a lot in
// system administration, e.g. when a user starts too many processes and
// therefore makes the system unresponsive for other users.
type ULimit struct {
	Soft int64 `yaml:"soft"`
	Hard int64 `yaml:"hard"`
}

// Compose the docker-compose.yml struct
type Compose struct {
	Version  string              `yaml:"version,omitempty"`
	Services map[string]*Service `yaml:"services,omitempty"`
}
