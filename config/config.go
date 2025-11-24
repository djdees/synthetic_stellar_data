package config

import (
	"fmt"
	"os"
	"strings"

  yaml "gopkg.in/yaml.v3"
)

// CassandraConfig holds Cassandra connection settings
// This structure maps directly to the YAML configuration file
type CassandraConfig struct {
	// Hosts is the list of Cassandra node addresses
	Hosts []string `yaml:"hosts"`

	// Keyspace is the name of the keyspace to use (will be created if it doesn't exist)
	Keyspace string `yaml:"keyspace"`

	// Consistency level for read/write operations
	// Valid values: ANY, ONE, TWO, THREE, QUORUM, ALL, LOCAL_QUORUM, EACH_QUORUM, LOCAL_ONE
	Consistency string `yaml:"consistency"`

	// ReplicationClass defines the replication strategy
	// Valid values: SimpleStrategy, NetworkTopologyStrategy
	ReplicationClass string `yaml:"replication_class"`

	// ReplicationFactor is the number of replicas for SimpleStrategy
	ReplicationFactor int `yaml:"replication_factor"`

	// DataCenters is used for NetworkTopologyStrategy (map of DC name to replication factor)
	DataCenters map[string]int `yaml:"data_centers,omitempty"`

	// Authentication credentials
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	// Optional: Connection settings
	Port             int  `yaml:"port,omitempty"`              // Default: 9042
	Timeout          int  `yaml:"timeout,omitempty"`           // Connection timeout in seconds
	ConnectTimeout   int  `yaml:"connect_timeout,omitempty"`   // Connect timeout in seconds
	NumConns         int  `yaml:"num_conns,omitempty"`         // Number of connections per host
	DisableInitialHostLookup bool `yaml:"disable_initial_host_lookup,omitempty"`
}

// LoadCassandraConfig loads Cassandra configuration from a YAML file
func LoadCassandraConfig(filename string) (*CassandraConfig, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var cfg CassandraConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// Set defaults for required fields
	setDefaults(&cfg)

	// Validate configuration
	if err := validateCassandraConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for optional configuration fields
func setDefaults(cfg *CassandraConfig) {
	// Default hosts
	if len(cfg.Hosts) == 0 {
		cfg.Hosts = []string{"127.0.0.1"}
	}

	// Default keyspace
	if cfg.Keyspace == "" {
		cfg.Keyspace = "stellargen"
	}

	// Default consistency level
	if cfg.Consistency == "" {
		cfg.Consistency = "QUORUM"
	}

	// Default replication strategy
	if cfg.ReplicationClass == "" {
		cfg.ReplicationClass = "SimpleStrategy"
	}

	// Default replication factor
	if cfg.ReplicationFactor == 0 && cfg.ReplicationClass == "SimpleStrategy" {
		cfg.ReplicationFactor = 1
	}

	// Default port
	if cfg.Port == 0 {
		cfg.Port = 9042
	}

	// Default timeout (30 seconds)
	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}

	// Default connect timeout (10 seconds)
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10
	}

	// Default number of connections per host
	if cfg.NumConns == 0 {
		cfg.NumConns = 2
	}
}

// validateCassandraConfig validates the Cassandra configuration
func validateCassandraConfig(cfg *CassandraConfig) error {
	// Validate hosts
	if len(cfg.Hosts) == 0 {
		return fmt.Errorf("at least one host must be specified")
	}

	for i, host := range cfg.Hosts {
		if strings.TrimSpace(host) == "" {
			return fmt.Errorf("host at index %d is empty", i)
		}
	}

	// Validate keyspace name
	if cfg.Keyspace == "" {
		return fmt.Errorf("keyspace name cannot be empty")
	}
	if !isValidKeyspaceName(cfg.Keyspace) {
		return fmt.Errorf("invalid keyspace name '%s': must contain only alphanumeric characters and underscores", cfg.Keyspace)
	}

	// Validate consistency level
	validConsistencies := map[string]bool{
		"ANY":         true,
		"ONE":         true,
		"TWO":         true,
		"THREE":       true,
		"QUORUM":      true,
		"ALL":         true,
		"LOCAL_QUORUM": true,
		"EACH_QUORUM": true,
		"LOCAL_ONE":   true,
	}
	if !validConsistencies[strings.ToUpper(cfg.Consistency)] {
		return fmt.Errorf("invalid consistency level '%s'", cfg.Consistency)
	}
	cfg.Consistency = strings.ToUpper(cfg.Consistency)

	// Validate replication strategy
	validStrategies := map[string]bool{
		"SimpleStrategy":            true,
		"NetworkTopologyStrategy":   true,
	}
	if !validStrategies[cfg.ReplicationClass] {
		return fmt.Errorf("invalid replication class '%s', must be SimpleStrategy or NetworkTopologyStrategy", cfg.ReplicationClass)
	}

	// Validate replication factor for SimpleStrategy
	if cfg.ReplicationClass == "SimpleStrategy" {
		if cfg.ReplicationFactor < 1 {
			return fmt.Errorf("replication_factor must be at least 1 for SimpleStrategy")
		}
		if cfg.ReplicationFactor > 3 {
			fmt.Fprintf(os.Stderr, "Warning: High replication factor (%d) may impact performance\n", cfg.ReplicationFactor)
		}
	}

	// Validate data centers for NetworkTopologyStrategy
	if cfg.ReplicationClass == "NetworkTopologyStrategy" {
		if len(cfg.DataCenters) == 0 {
			return fmt.Errorf("data_centers must be specified for NetworkTopologyStrategy")
		}
		for dc, rf := range cfg.DataCenters {
			if rf < 1 {
				return fmt.Errorf("replication factor for datacenter '%s' must be at least 1", dc)
			}
		}
	}

	// Validate port
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port %d, must be between 1 and 65535", cfg.Port)
	}

	// Validate timeouts
	if cfg.Timeout < 1 {
		return fmt.Errorf("timeout must be at least 1 second")
	}
	if cfg.ConnectTimeout < 1 {
		return fmt.Errorf("connect_timeout must be at least 1 second")
	}

	// Validate connection count
	if cfg.NumConns < 1 {
		return fmt.Errorf("num_conns must be at least 1")
	}
	if cfg.NumConns > 10 {
		fmt.Fprintf(os.Stderr, "Warning: High connection count (%d) per host\n", cfg.NumConns)
	}

	return nil
}

// isValidKeyspaceName checks if a keyspace name is valid
// Cassandra keyspace names must be alphanumeric or underscore, and not start with a number
func isValidKeyspaceName(name string) bool {
	if len(name) == 0 || len(name) > 48 {
		return false
	}

	for i, r := range name {
		if i == 0 {
			// First character cannot be a number
			if r >= '0' && r <= '9' {
				return false
			}
		}
		// Must be alphanumeric or underscore
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}

	return true
}

// GetReplicationString returns the replication string for CREATE KEYSPACE
func (cfg *CassandraConfig) GetReplicationString() string {
	if cfg.ReplicationClass == "SimpleStrategy" {
		return fmt.Sprintf("{'class': 'SimpleStrategy', 'replication_factor': %d}", cfg.ReplicationFactor)
	}

	// NetworkTopologyStrategy
	var parts []string
	parts = append(parts, "'class': 'NetworkTopologyStrategy'")
	for dc, rf := range cfg.DataCenters {
		parts = append(parts, fmt.Sprintf("'%s': %d", dc, rf))
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

// PrintConfig prints a formatted configuration summary
func (cfg *CassandraConfig) PrintConfig() {
	fmt.Println("Cassandra Configuration:")
	fmt.Printf("  Hosts:              %s\n", strings.Join(cfg.Hosts, ", "))
	fmt.Printf("  Port:               %d\n", cfg.Port)
	fmt.Printf("  Keyspace:           %s\n", cfg.Keyspace)
	fmt.Printf("  Consistency:        %s\n", cfg.Consistency)
	fmt.Printf("  Replication Class:  %s\n", cfg.ReplicationClass)
	
	if cfg.ReplicationClass == "SimpleStrategy" {
		fmt.Printf("  Replication Factor: %d\n", cfg.ReplicationFactor)
	} else {
		fmt.Printf("  Data Centers:\n")
		for dc, rf := range cfg.DataCenters {
			fmt.Printf("    %s: %d\n", dc, rf)
		}
	}
	
	if cfg.Username != "" {
		fmt.Printf("  Authentication:     Enabled (user: %s)\n", cfg.Username)
	} else {
		fmt.Printf("  Authentication:     Disabled\n")
	}
	
	fmt.Printf("  Timeout:            %ds\n", cfg.Timeout)
	fmt.Printf("  Connect Timeout:    %ds\n", cfg.ConnectTimeout)
	fmt.Printf("  Connections/Host:   %d\n", cfg.NumConns)
	fmt.Println()
}

// SaveConfig saves the configuration to a YAML file
// Useful for generating template configurations
func SaveConfig(cfg *CassandraConfig, filename string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// DefaultConfig returns a CassandraConfig with default values
// Useful for generating example configurations
func DefaultConfig() *CassandraConfig {
	return &CassandraConfig{
		Hosts:             []string{"127.0.0.1"},
		Port:              9042,
		Keyspace:          "stellargen",
		Consistency:       "QUORUM",
		ReplicationClass:  "SimpleStrategy",
		ReplicationFactor: 1,
		Username:          "",
		Password:          "",
		Timeout:           30,
		ConnectTimeout:    10,
		NumConns:          2,
	}
}

// MultiDCConfig returns an example NetworkTopologyStrategy configuration
func MultiDCConfig() *CassandraConfig {
	return &CassandraConfig{
		Hosts:            []string{"10.0.1.1", "10.0.2.1", "10.0.3.1"},
		Port:             9042,
		Keyspace:         "stellargen",
		Consistency:      "LOCAL_QUORUM",
		ReplicationClass: "NetworkTopologyStrategy",
		DataCenters: map[string]int{
			"dc1": 3,
			"dc2": 2,
		},
		Username:       "stellargen_user",
		Password:       "secure_password",
		Timeout:        30,
		ConnectTimeout: 10,
		NumConns:       3,
	}
}
