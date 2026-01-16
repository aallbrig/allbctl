package cmd

import (
	"testing"
)

func TestExtractCloudCLIVersion(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		output   string
		expected string
	}{
		{
			name:     "AWS CLI version",
			cli:      "aws",
			output:   "aws-cli/2.33.1 Python/3.13.11 Linux/6.14.0-37-generic exe/x86_64.linuxmint.22",
			expected: "2.33.1",
		},
		{
			name:     "kubectl version short",
			cli:      "kubectl",
			output:   "v1.28.0",
			expected: "1.28.0",
		},
		{
			name:     "kubectl version long",
			cli:      "kubectl",
			output:   "Client Version: v1.28.0",
			expected: "1.28.0",
		},
		{
			name:     "gcloud version",
			cli:      "gcloud",
			output:   "478.0.0",
			expected: "478.0.0",
		},
		{
			name:     "azure CLI version",
			cli:      "az",
			output:   "2.56.0",
			expected: "2.56.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCloudCLIVersion(tt.cli, tt.output)
			if result != tt.expected {
				t.Errorf("extractCloudCLIVersion(%q, %q) = %q; want %q", tt.cli, tt.output, result, tt.expected)
			}
		})
	}
}

func TestDetectCloudCLIs(t *testing.T) {
	// This test verifies the function runs without errors
	// Actual CLI detection depends on system state
	clis := detectCloudCLIs()

	// Should return a slice (may be empty if no CLIs installed)
	if clis == nil {
		t.Error("detectCloudCLIs() returned nil; want non-nil slice")
	}
}

func TestGetAWSProfiles(t *testing.T) {
	// This test verifies the function runs without errors
	// Actual profiles depend on system configuration
	profiles := getAWSProfiles()

	// Should return a slice (may be empty if no profiles configured)
	if profiles == nil {
		t.Error("getAWSProfiles() returned nil; want non-nil slice")
	}
}

func TestGetGCloudProfiles(t *testing.T) {
	// This test verifies the function runs without errors
	// Actual profiles depend on system configuration
	profiles := getGCloudProfiles()

	// Should return a slice (may be empty if gcloud not installed)
	if profiles == nil {
		t.Error("getGCloudProfiles() returned nil; want non-nil slice")
	}
}

func TestGetAzureProfiles(t *testing.T) {
	// This test verifies the function runs without errors
	// Actual profiles depend on system configuration
	profiles := getAzureProfiles()

	// Should return a slice (may be empty if az not installed)
	if profiles == nil {
		t.Error("getAzureProfiles() returned nil; want non-nil slice")
	}
}

func TestGetKubeContexts(t *testing.T) {
	// This test verifies the function runs without errors
	// Actual contexts depend on system configuration
	contexts := getKubeContexts()

	// Should return a slice (may be empty if kubectl not installed)
	if contexts == nil {
		t.Error("getKubeContexts() returned nil; want non-nil slice")
	}
}

func TestGetAWSRegions(t *testing.T) {
	// This test verifies the function returns a non-empty slice
	// We can't test the actual API call without AWS credentials,
	// but we document that the function should return all AWS regions
	t.Run("returns comprehensive region list", func(t *testing.T) {
		// Function should return all known AWS regions when no specific region is configured
		// The actual test of the function would require AWS credentials
		t.Skip("getAWSRegions requires AWS credentials to test API call; fallback list is verified in implementation")
	})
}

func TestExtractResourceTypeFromARN(t *testing.T) {
	tests := []struct {
		name     string
		arn      string
		expected string
	}{
		{
			name:     "EC2 instance",
			arn:      "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
			expected: "ec2::instance",
		},
		{
			name:     "S3 bucket",
			arn:      "arn:aws:s3:::my-bucket",
			expected: "s3::my-bucket",
		},
		{
			name:     "Lambda function",
			arn:      "arn:aws:lambda:us-west-2:123456789012:function/my-function",
			expected: "lambda::function",
		},
		{
			name:     "GameLift fleet",
			arn:      "arn:aws:gamelift:us-east-1:123456789012:fleet/fleet-12345678-1234-1234-1234-123456789012",
			expected: "gamelift::fleet",
		},
		{
			name:     "RDS instance",
			arn:      "arn:aws:rds:us-east-1:123456789012:db/mydb",
			expected: "rds::db",
		},
		{
			name:     "DynamoDB table",
			arn:      "arn:aws:dynamodb:us-east-1:123456789012:table/my-table",
			expected: "dynamodb::table",
		},
		{
			name:     "Invalid ARN",
			arn:      "invalid-arn",
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractResourceTypeFromARN(tt.arn)
			if result != tt.expected {
				t.Errorf("extractResourceTypeFromARN(%q) = %q; want %q", tt.arn, result, tt.expected)
			}
		})
	}
}

func TestPrintAWSDetailsWithProfiles(t *testing.T) {
	// This test documents the expected behavior with profile filtering
	// Actual test would require AWS credentials and profiles
	t.Run("filters to specific profile when --profile flag is set", func(t *testing.T) {
		t.Skip("Requires AWS credentials and multiple profiles to test profile filtering")
	})

	t.Run("processes all profiles when no --profile flag is set", func(t *testing.T) {
		t.Skip("Requires AWS credentials and multiple profiles to test profile iteration")
	})

	t.Run("shows error when specified profile does not exist", func(t *testing.T) {
		t.Skip("Requires AWS credentials to test profile validation")
	})
}
