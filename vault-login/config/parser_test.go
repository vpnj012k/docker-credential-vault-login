// Tests to write:
// * Order of file reading 

package config

import (
	"fmt"
	"os"
	"testing"
	"github.com/mitchellh/go-homedir"
	test "github.com/morningconsult/docker-credential-vault-login/vault-login/testing"
)

// TestReadsFileEnv tests that the GetCredHelperConfig function
// looks for and parses the config file specified by the 
// DOCKER_CREDS_CONFIG_FILE environment variable
func TestReadsFileEnv(t *testing.T) {
        const testFilePath = "/tmp/docker-credential-vault-login-testfile.json"
        cfg := &CredHelperConfig{
                Method:   VaultAuthMethodAWSIAM,
                Role:     "dev-role-iam",
                Secret:   "secret/foo/bar",
                ServerID: "vault.example.com",
        }
        data := test.EncodeJSON(t, cfg)
        test.MakeFile(t, testFilePath, data)
        defer test.DeleteFile(t, testFilePath)

        path := os.Getenv(EnvConfigFilePath)
        os.Setenv(EnvConfigFilePath, testFilePath)
        defer os.Setenv(EnvConfigFilePath, path)

        if _, err := GetCredHelperConfig(); err != nil {
                t.Errorf("Failed to read config file specified by environment variable %s", 
                        EnvConfigFilePath)
        }
}

// TestConfigFileMissing tests that if the config file located
// at either the default path or the path given by the 
// DOCKER_CREDS_CONFIG_FILE environment variable does not exist,
// GetCredHelperConfig() throws the appropriate error
func TestConfigFileMissing(t *testing.T) {
        const testFilePath = "/tmp/docker-credential-vault-login-testfile-2.json"

        path := os.Getenv(EnvConfigFilePath)
        os.Setenv(EnvConfigFilePath, testFilePath)
        defer os.Setenv(EnvConfigFilePath, path)

        if _, err := GetCredHelperConfig(); err != nil {
                if !os.IsNotExist(err) {
                        t.Errorf("%s (expected os.ErrNotExist, got %v)", 
                                "GetCredHelperConfig() returned unexpected error", 
                                err)
                }
        } else {
                t.Fatal("Expected to receive an error but didn't")
        }
}

// TestEmptyConfigFile tests that if the configuration file is
// just an empty JSON, the expected errors are returned.
func TestEmptyConfigFile(t *testing.T) {
        const testFilePath = "/tmp/docker-credential-vault-login-testfile-3.json"
        var expectedError = fmt.Sprintf("%s\n%s\n%s",
                fmt.Sprintf("Configuration file %s has the following errors:", testFilePath),
                "* No Vault authentication method (\"vault_auth_method\") is provided",
                "* No path to the location of your secret in Vault (\"vault_secret_path\") is provided")

        test.MakeFile(t, testFilePath, []byte("{}"))
        defer test.DeleteFile(t, testFilePath)

        path := os.Getenv(EnvConfigFilePath)
        os.Setenv(EnvConfigFilePath, testFilePath)
        defer os.Setenv(EnvConfigFilePath, path)

        if _, err := GetCredHelperConfig(); err != nil {
                test.ErrorsEqual(t, err, expectedError)
        } else {
                t.Fatal("Expected to receive an error but didn't")
        }
}

// TestConfigMissingMethod tests that GetCredHelperConfig 
// return the expected error message when no authentication
// method is provided in the configuration file.
func TestConfigMissingMethod(t *testing.T) {
        const testFilePath = "/tmp/docker-credential-vault-login-testfile-4.json"
        var expectedError = fmt.Sprintf("%s\n%s",
                fmt.Sprintf("Configuration file %s has the following errors:", testFilePath),
                "* No Vault authentication method (\"vault_auth_method\") is provided")
        
        cfg := &CredHelperConfig{
                Role:     "dev-role-iam",
                Secret:   "secret/foo/bar",
                ServerID: "vault.example.com",
        }
        data := test.EncodeJSON(t, cfg)
        test.MakeFile(t, testFilePath, data)
        defer test.DeleteFile(t, testFilePath)

        path := os.Getenv(EnvConfigFilePath)
        os.Setenv(EnvConfigFilePath, testFilePath)
        defer os.Setenv(EnvConfigFilePath, path)

        if _, err := GetCredHelperConfig(); err != nil {
                test.ErrorsEqual(t, err, expectedError)
        } else {
                t.Fatal("Expected to receive an error but didn't")
        }
}

// TestConfigMissingSecret tests that GetCredHelperConfig 
// returns the expected error message when no path to a Vault
// secret is provided in the configuration file. This secret 
// is the location in Vault at which your Docker credentials 
// are stored
func TestConfigMissingSecret(t *testing.T) {
        const testFilePath = "/tmp/docker-credential-vault-login-testfile-5.json"
        var expectedError = fmt.Sprintf("%s\n%s",
                fmt.Sprintf("Configuration file %s has the following errors:", testFilePath),
                "* No path to the location of your secret in Vault (\"vault_secret_path\") is provided")
        
        cfg := &CredHelperConfig{
                Method:   VaultAuthMethodAWSIAM,
                Role:     "dev-role-iam",
                ServerID: "vault.example.com",
        }
        data := test.EncodeJSON(t, cfg)
        test.MakeFile(t, testFilePath, data)
        defer test.DeleteFile(t, testFilePath)

        path := os.Getenv(EnvConfigFilePath)
        os.Setenv(EnvConfigFilePath, testFilePath)
        defer os.Setenv(EnvConfigFilePath, path)

        if _, err := GetCredHelperConfig(); err != nil {
                test.ErrorsEqual(t, err, expectedError)
        } else {
                t.Fatal("Expected to receive an error but didn't")
        }
}

// TestConfigMissingToken tests that GetCredHelperConfig
// returns the expected error message when "token" is selected
// as the Vault authentication method in the config file but
// the VAULT_TOKEN environment variable is not set
func TestConfigMissingToken(t *testing.T) {
        const testFilePath = "/tmp/docker-credential-vault-login-testfile-6.json"
        var expectedError = fmt.Sprintf("%s\n%s",
                fmt.Sprintf("Configuration file %s has the following errors:", testFilePath),
                "* VAULT_TOKEN environment variable is not set")
        
        cfg := &CredHelperConfig{
                Method:   VaultAuthMethodToken,
                Secret:   "secret/foo/bar",
        }
        data := test.EncodeJSON(t, cfg)
        test.MakeFile(t, testFilePath, data)
        defer test.DeleteFile(t, testFilePath)

        path := os.Getenv(EnvConfigFilePath)
        os.Setenv(EnvConfigFilePath, testFilePath)
        defer os.Setenv(EnvConfigFilePath, path)

        originalToken := os.Getenv("VAULT_TOKEN")
        defer os.Setenv("VAULT_TOKEN", originalToken)

        os.Setenv("VAULT_TOKEN", "")
        if _, err := GetCredHelperConfig(); err != nil {
                test.ErrorsEqual(t, err, expectedError)
        } else {
                t.Fatal("Expected to receive an error but didn't")
        }
}

// TestConfigMissingRole tests that GetCredHelperConfig 
// returns the expected error message when no Vault role
// is provided in the configuration file when the AWS auth
// method is chosen
func TestConfigMissingRole(t *testing.T) {
        const testFilePath = "/tmp/docker-credential-vault-login-testfile-7.json"
        var expectedError = fmt.Sprintf("%s\n%s%s",
                fmt.Sprintf("Configuration file %s has the following errors:", testFilePath),
                "* No Vault role (\"vault_role\") is provided (required when ",
                "the AWS authentication method is chosen)")
        
        cfg := &CredHelperConfig{
                Method:   VaultAuthMethodAWSIAM,
                Secret:   "secret/foo/bar",
        }
        data := test.EncodeJSON(t, cfg)
        test.MakeFile(t, testFilePath, data)
        defer test.DeleteFile(t, testFilePath)

        path := os.Getenv(EnvConfigFilePath)
        os.Setenv(EnvConfigFilePath, testFilePath)
        defer os.Setenv(EnvConfigFilePath, path)

        if _, err := GetCredHelperConfig(); err != nil {
                test.ErrorsEqual(t, err, expectedError)
        } else {
                t.Fatal("Expected to receive an error but didn't")
        }

}

// TestConfigBadAuthMethod tests that if an unsupported
// authentication method is provided in the "vault_auth_method"
// field of the config.json file, GetCredHelperConfig returns
// the appropriate error.
func TestConfigBadAuthMethod(t *testing.T) {
        const (
                testFilePath = "/tmp/docker-credential-vault-login-testfile-8.json"
                badMethod = "potato"
        )
        
        var expectedError = fmt.Sprintf("%s\n* %s",
                fmt.Sprintf("Configuration file %s has the following errors:", testFilePath),
                fmt.Sprintf(`Unrecognized Vault authentication method ("vault_auth_method") value %q (must be one of "iam", "ec2", or "token")`, badMethod),
        )
        
        cfg := &CredHelperConfig{
                Method:   VaultAuthMethod("potato"),
                Secret:   "secret/foo/bar",
        }
        data := test.EncodeJSON(t, cfg)
        test.MakeFile(t, testFilePath, data)
        defer test.DeleteFile(t, testFilePath)

        path := os.Getenv(EnvConfigFilePath)
        os.Setenv(EnvConfigFilePath, testFilePath)
        defer os.Setenv(EnvConfigFilePath, path)

        if _, err := GetCredHelperConfig(); err != nil {
                test.ErrorsEqual(t, err, expectedError)
        } else {
                t.Fatal("Expected to receive an error but didn't")
        }
}

// TestConfigAllowsTildeInPath tests that  GetCredHelperConfig
// does not fail when the path to the config.json file includes
// a tile (e.g. "~/Desktop/config.json") in it.
func TestConfigAllowsTildeInPath(t *testing.T) {
	const testFilePath = "~/docker-credential-vault-login-testfile-9.json"
	cfg := &CredHelperConfig{
		Method:   VaultAuthMethodAWSIAM,
		Role:     "dev-role-iam",
		Secret:   "secret/foo/bar",
		ServerID: "vault.example.com",
	}
	
	pathExpanded, err := homedir.Expand(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	data := test.EncodeJSON(t, cfg)
	test.MakeFile(t, pathExpanded, data)
	defer test.DeleteFile(t, pathExpanded)

	path := os.Getenv(EnvConfigFilePath)
	// Set $DOCKER_CREDS_CONFIG_FILE environment variable to testFilePath
	os.Setenv(EnvConfigFilePath, testFilePath)
	defer os.Setenv(EnvConfigFilePath, path)

	if _, err := GetCredHelperConfig(); err != nil {
		t.Errorf("Failed to read config file specified by environment variable %s", 
			EnvConfigFilePath)
	}
}