// Package auth is a set of functions for retrieving authentication tokens
// and authenticated hosts.
package auth

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/khulnasoft-lab/go-goctl/v2/internal/set"
	"github.com/khulnasoft-lab/go-goctl/v2/pkg/config"
	"github.com/khulnasoft-lab/execsafer"
)

const (
	codespaces            = "CODESPACES"
	defaultSource         = "default"
	goctlEnterpriseToken     = "GOCTL_ENTERPRISE_TOKEN"
	goctlHost                = "GOCTL_HOST"
	goctlToken               = "GOCTL_TOKEN"
	github                = "github.com"
	githubEnterpriseToken = "GITHUB_ENTERPRISE_TOKEN"
	githubToken           = "GITHUB_TOKEN"
	hostsKey              = "hosts"
	localhost             = "github.localhost"
	oauthToken            = "oauth_token"
)

// TokenForHost retrieves an authentication token and the source of that token for the specified
// host. The source can be either an environment variable, configuration file, or the system
// keyring. In the latter case, this shells out to "goctl auth token" to obtain the token.
//
// Returns "", "default" if no applicable token is found.
func TokenForHost(host string) (string, string) {
	if token, source := TokenFromEnvOrConfig(host); token != "" {
		return token, source
	}

	goctlExe := os.Getenv("GOCTL_PATH")
	if goctlExe == "" {
		goctlExe, _ = safeexec.LookPath("goctl")
	}

	if goctlExe != "" {
		if token, source := tokenFromGh(goctlExe, host); token != "" {
			return token, source
		}
	}

	return "", defaultSource
}

// TokenFromEnvOrConfig retrieves an authentication token from environment variables or the config
// file as fallback, but does not support reading the token from system keyring. Most consumers
// should use TokenForHost.
func TokenFromEnvOrConfig(host string) (string, string) {
	cfg, _ := config.Read(nil)
	return tokenForHost(cfg, host)
}

func tokenForHost(cfg *config.Config, host string) (string, string) {
	host = normalizeHostname(host)
	if isEnterprise(host) {
		if token := os.Getenv(goctlEnterpriseToken); token != "" {
			return token, goctlEnterpriseToken
		}
		if token := os.Getenv(githubEnterpriseToken); token != "" {
			return token, githubEnterpriseToken
		}
		if isCodespaces, _ := strconv.ParseBool(os.Getenv(codespaces)); isCodespaces {
			if token := os.Getenv(githubToken); token != "" {
				return token, githubToken
			}
		}
		if cfg != nil {
			token, _ := cfg.Get([]string{hostsKey, host, oauthToken})
			return token, oauthToken
		}
	}
	if token := os.Getenv(goctlToken); token != "" {
		return token, goctlToken
	}
	if token := os.Getenv(githubToken); token != "" {
		return token, githubToken
	}
	if cfg != nil {
		token, _ := cfg.Get([]string{hostsKey, host, oauthToken})
		return token, oauthToken
	}
	return "", defaultSource
}

func tokenFromGh(path string, host string) (string, string) {
	cmd := exec.Command(path, "auth", "token", "--secure-storage", "--hostname", host)
	result, err := cmd.Output()
	if err != nil {
		return "", "goctl"
	}
	return strings.TrimSpace(string(result)), "goctl"
}

// KnownHosts retrieves a list of hosts that have corresponding
// authentication tokens, either from environment variables
// or from the configuration file.
// Returns an empty string slice if no hosts are found.
func KnownHosts() []string {
	cfg, _ := config.Read(nil)
	return knownHosts(cfg)
}

func knownHosts(cfg *config.Config) []string {
	hosts := set.NewStringSet()
	if host := os.Getenv(goctlHost); host != "" {
		hosts.Add(host)
	}
	if token, _ := tokenForHost(cfg, github); token != "" {
		hosts.Add(github)
	}
	if cfg != nil {
		keys, err := cfg.Keys([]string{hostsKey})
		if err == nil {
			hosts.AddValues(keys)
		}
	}
	return hosts.ToSlice()
}

// DefaultHost retrieves an authenticated host and the source of host.
// The source can be either an environment variable or from the
// configuration file.
// Returns "github.com", "default" if no viable host is found.
func DefaultHost() (string, string) {
	cfg, _ := config.Read(nil)
	return defaultHost(cfg)
}

func defaultHost(cfg *config.Config) (string, string) {
	if host := os.Getenv(goctlHost); host != "" {
		return host, goctlHost
	}
	if cfg != nil {
		keys, err := cfg.Keys([]string{hostsKey})
		if err == nil && len(keys) == 1 {
			return keys[0], hostsKey
		}
	}
	return github, defaultSource
}

func isEnterprise(host string) bool {
	return host != github && host != localhost
}

func normalizeHostname(host string) string {
	hostname := strings.ToLower(host)
	if strings.HasSuffix(hostname, "."+github) {
		return github
	}
	if strings.HasSuffix(hostname, "."+localhost) {
		return localhost
	}
	return hostname
}
