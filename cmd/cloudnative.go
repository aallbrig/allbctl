package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/spf13/cobra"
)

var (
	regionFlag  string
	profileFlag string
)

// CloudNativeCmd represents the cloud-native command
var CloudNativeCmd = &cobra.Command{
	Use:     "cloud-native",
	Aliases: []string{"cloud_native", "cn"},
	Short:   "Display cloud CLI information and resources",
	Long: `Display cloud CLI versions and resource information for AWS, GCP, Azure, and Kubernetes.

This command shows:
- CLI version for each detected cloud provider
- Number of configured profiles/users
- Resource counts by region (when detailed subcommands are used)

Examples:
  allbctl status cloud-native          # Show summary for all detected cloud CLIs
  allbctl status cn                    # Short alias
  allbctl status cloud-native aws      # Show detailed AWS resource info
  allbctl status cn aws --region us-east-1  # AWS resources in specific region`,
	Run: func(cmd *cobra.Command, args []string) {
		printCloudNativeSummary()
	},
}

// CloudNativeAWSCmd represents the AWS-specific command
var CloudNativeAWSCmd = &cobra.Command{
	Use:   "aws",
	Short: "Display detailed AWS resource information",
	Long: `Display detailed AWS resource information including:
- CLI version
- Configured profiles
- Resource counts by type and region (only regions with resources are shown)
- Uses AWS Resource Groups Tagging API to discover all resources

Examples:
  allbctl status cloud-native aws                           # All profiles, all regions with resources
  allbctl status cloud-native aws --region us-east-1        # All profiles, specific region
  allbctl status cloud-native aws --profile production      # Specific profile, all regions
  allbctl status cloud-native aws --profile prod --region us-east-1  # Specific profile and region`,
	Run: func(cmd *cobra.Command, args []string) {
		printAWSDetails()
	},
}

// CloudNativeGCPCmd represents the GCP-specific command
var CloudNativeGCPCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Display detailed GCP resource information",
	Long:  `Display detailed GCP resource information (implementation pending).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GCP detailed view: implementation todo")
	},
}

// CloudNativeAzureCmd represents the Azure-specific command
var CloudNativeAzureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Display detailed Azure resource information",
	Long:  `Display detailed Azure resource information (implementation pending).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Azure detailed view: implementation todo")
	},
}

// CloudNativeK8sCmd represents the Kubernetes-specific command
var CloudNativeK8sCmd = &cobra.Command{
	Use:     "kubernetes",
	Aliases: []string{"k8s"},
	Short:   "Display detailed Kubernetes resource information",
	Long:    `Display detailed Kubernetes resource information (implementation pending).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kubernetes detailed view: implementation todo")
	},
}

func init() {
	// Add subcommands to cloud-native
	CloudNativeCmd.AddCommand(CloudNativeAWSCmd)
	CloudNativeCmd.AddCommand(CloudNativeGCPCmd)
	CloudNativeCmd.AddCommand(CloudNativeAzureCmd)
	CloudNativeCmd.AddCommand(CloudNativeK8sCmd)

	// Add flags
	CloudNativeAWSCmd.Flags().StringVar(&regionFlag, "region", "", "Specific AWS region to query")
	CloudNativeAWSCmd.Flags().StringVar(&profileFlag, "profile", "", "Specific AWS profile to query")
}

// CloudCLIInfo holds information about a cloud CLI
type CloudCLIInfo struct {
	Name             string
	Version          string
	KustomizeVersion string // For kubectl
	ProfileCount     int
	Profiles         []string
	Connected        bool
}

// checkAWSConnectivity checks if AWS CLI can connect to AWS (any profile connected = true)
func checkAWSConnectivity() bool {
	profiles := getAWSProfiles()
	if len(profiles) == 0 {
		return false
	}

	// Check each profile concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	anyConnected := false

	for _, profile := range profiles {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if checkAWSProfileConnectivity(p) {
				mu.Lock()
				anyConnected = true
				mu.Unlock()
			}
		}(profile)
	}
	wg.Wait()

	return anyConnected
}

// checkAWSProfileConnectivity checks if a specific AWS profile can connect
func checkAWSProfileConnectivity(profile string) bool {
	cmd := exec.Command("aws", "sts", "get-caller-identity", "--profile", profile)
	err := cmd.Run()
	return err == nil
}

// checkGCloudConnectivity checks if gcloud CLI can connect to GCP
func checkGCloudConnectivity() bool {
	cmd := exec.Command("gcloud", "auth", "list", "--filter=status:ACTIVE", "--format=value(account)")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

// checkAzureConnectivity checks if Azure CLI can connect to Azure
func checkAzureConnectivity() bool {
	cmd := exec.Command("az", "account", "show")
	err := cmd.Run()
	return err == nil
}

// checkKubectlConnectivity checks if kubectl can connect to a cluster
func checkKubectlConnectivity() bool {
	cmd := exec.Command("kubectl", "cluster-info")
	err := cmd.Run()
	return err == nil
}

// getKustomizeVersion gets the kustomize version from kubectl
func getKustomizeVersion() string {
	cmd := exec.Command("kubectl", "version", "--client", "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Parse JSON to extract kustomizeVersion
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return ""
	}

	// kustomizeVersion is at the top level of the JSON
	if kustomizeVersion, ok := result["kustomizeVersion"].(string); ok {
		return strings.TrimPrefix(kustomizeVersion, "v")
	}

	return ""
}

// detectCloudCLIs detects installed cloud CLIs and their info
func detectCloudCLIs() []CloudCLIInfo {
	var clis []CloudCLIInfo
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Check AWS CLI
	wg.Add(1)
	go func() {
		defer wg.Done()
		if exists("aws") {
			info := CloudCLIInfo{Name: "aws"}
			if version := getCloudCLIVersion("aws"); version != "" {
				info.Version = version
			}
			if profiles := getAWSProfiles(); len(profiles) >= 0 {
				info.ProfileCount = len(profiles)
				info.Profiles = profiles
			}
			info.Connected = checkAWSConnectivity()
			mu.Lock()
			clis = append(clis, info)
			mu.Unlock()
		}
	}()

	// Check gcloud CLI
	wg.Add(1)
	go func() {
		defer wg.Done()
		if exists("gcloud") {
			info := CloudCLIInfo{Name: "gcloud"}
			if version := getCloudCLIVersion("gcloud"); version != "" {
				info.Version = version
			}
			if profiles := getGCloudProfiles(); len(profiles) >= 0 {
				info.ProfileCount = len(profiles)
				info.Profiles = profiles
			}
			info.Connected = checkGCloudConnectivity()
			mu.Lock()
			clis = append(clis, info)
			mu.Unlock()
		}
	}()

	// Check Azure CLI
	wg.Add(1)
	go func() {
		defer wg.Done()
		if exists("az") {
			info := CloudCLIInfo{Name: "az"}
			if version := getCloudCLIVersion("az"); version != "" {
				info.Version = version
			}
			if profiles := getAzureProfiles(); len(profiles) >= 0 {
				info.ProfileCount = len(profiles)
				info.Profiles = profiles
			}
			info.Connected = checkAzureConnectivity()
			mu.Lock()
			clis = append(clis, info)
			mu.Unlock()
		}
	}()

	// Check kubectl
	wg.Add(1)
	go func() {
		defer wg.Done()
		if exists("kubectl") {
			info := CloudCLIInfo{Name: "kubectl"}
			if version := getCloudCLIVersion("kubectl"); version != "" {
				info.Version = version
			}
			if kustomizeVersion := getKustomizeVersion(); kustomizeVersion != "" {
				info.KustomizeVersion = kustomizeVersion
			}
			if contexts := getKubeContexts(); len(contexts) >= 0 {
				info.ProfileCount = len(contexts)
				info.Profiles = contexts
			}
			info.Connected = checkKubectlConnectivity()
			mu.Lock()
			clis = append(clis, info)
			mu.Unlock()
		}
	}()

	wg.Wait()
	return clis
}

// getCloudCLIVersion gets the version of a cloud CLI
func getCloudCLIVersion(cli string) string {
	var cmd *exec.Cmd

	switch cli {
	case "aws":
		cmd = exec.Command("aws", "--version")
	case "gcloud":
		cmd = exec.Command("gcloud", "version", "--format", "value(version)")
	case "az":
		cmd = exec.Command("az", "version", "--query", "\"azure-cli\"", "-o", "tsv")
	case "kubectl":
		cmd = exec.Command("kubectl", "version", "--client")
	default:
		return ""
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	return extractCloudCLIVersion(cli, version)
}

// extractCloudCLIVersion extracts clean version from cloud CLI output
func extractCloudCLIVersion(cli, output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}

	lines := strings.Split(output, "\n")
	firstLine := strings.TrimSpace(lines[0])

	switch cli {
	case "aws":
		// "aws-cli/2.33.1 Python/3.13.11 Linux/6.14.0-37-generic exe/x86_64.linuxmint.22"
		if strings.HasPrefix(firstLine, "aws-cli/") {
			parts := strings.Fields(firstLine)
			if len(parts) >= 1 {
				return strings.TrimPrefix(parts[0], "aws-cli/")
			}
		}
	case "gcloud":
		// Just the version number with --format value(version)
		return firstLine
	case "az":
		// Just the version number with --query
		return firstLine
	case "kubectl":
		// "Client Version: v1.28.0" or "v1.28.0"
		if strings.Contains(firstLine, "Client Version:") {
			parts := strings.Fields(firstLine)
			for _, part := range parts {
				if strings.HasPrefix(part, "v") {
					return strings.TrimPrefix(part, "v")
				}
			}
		}
		return strings.TrimPrefix(firstLine, "v")
	}

	return firstLine
}

// getAWSProfiles returns list of AWS profiles
func getAWSProfiles() []string {
	cmd := exec.Command("aws", "configure", "list-profiles")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	profilesStr := strings.TrimSpace(string(output))
	if profilesStr == "" {
		return []string{}
	}

	profiles := strings.Split(profilesStr, "\n")
	var result []string
	for _, p := range profiles {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// getDefaultAWSProfile returns the default AWS profile name
func getDefaultAWSProfile() string {
	// Check AWS_PROFILE environment variable first
	cmd := exec.Command("sh", "-c", "echo $AWS_PROFILE")
	output, err := cmd.Output()
	if err == nil {
		profile := strings.TrimSpace(string(output))
		if profile != "" {
			return profile
		}
	}

	// If no AWS_PROFILE set, default profile is "default"
	return "default"
}

// getDefaultRegionForProfile returns the default region for a specific AWS profile
func getDefaultRegionForProfile(profile string) string {
	cmd := exec.Command("aws", "configure", "get", "region", "--profile", profile)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getGCloudProfiles returns list of GCP accounts
func getGCloudProfiles() []string {
	cmd := exec.Command("gcloud", "auth", "list", "--format", "value(account)")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	accountsStr := strings.TrimSpace(string(output))
	if accountsStr == "" {
		return []string{}
	}

	accounts := strings.Split(accountsStr, "\n")
	var result []string
	for _, a := range accounts {
		a = strings.TrimSpace(a)
		if a != "" {
			result = append(result, a)
		}
	}
	return result
}

// getAzureProfiles returns list of Azure accounts
func getAzureProfiles() []string {
	cmd := exec.Command("az", "account", "list", "--query", "[].user.name", "-o", "tsv")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	accountsStr := strings.TrimSpace(string(output))
	if accountsStr == "" {
		return []string{}
	}

	accounts := strings.Split(accountsStr, "\n")
	var result []string
	for _, a := range accounts {
		a = strings.TrimSpace(a)
		if a != "" {
			result = append(result, a)
		}
	}
	return result
}

// getKubeContexts returns list of kubectl contexts
func getKubeContexts() []string {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	contextsStr := strings.TrimSpace(string(output))
	if contextsStr == "" {
		return []string{}
	}

	contexts := strings.Split(contextsStr, "\n")
	var result []string
	for _, c := range contexts {
		c = strings.TrimSpace(c)
		if c != "" {
			result = append(result, c)
		}
	}
	return result
}

// printCloudNativeSummary prints a summary of all cloud CLIs
func printCloudNativeSummary() {
	clis := detectCloudCLIs()

	if len(clis) == 0 {
		// No output if no CLIs detected
		return
	}

	fmt.Println("Cloud Native:")
	for _, cli := range clis {
		// Print CLI name and version
		if cli.Version != "" {
			fmt.Printf("  %s (%s)", cli.Name, cli.Version)
		} else {
			fmt.Printf("  %s", cli.Name)
		}

		// For kubectl, also show kustomize version
		if cli.Name == "kubectl" && cli.KustomizeVersion != "" {
			fmt.Printf(" [kustomize: %s]", cli.KustomizeVersion)
		}

		// Print profile/context count (kubectl uses contexts, others use profiles)
		if cli.Name == "kubectl" {
			if cli.ProfileCount == 0 {
				fmt.Printf(" - 0 contexts")
			} else if cli.ProfileCount == 1 {
				fmt.Printf(" - 1 context")
			} else {
				fmt.Printf(" - %d contexts", cli.ProfileCount)
			}
		} else {
			if cli.ProfileCount == 0 {
				fmt.Printf(" - 0 profiles")
			} else if cli.ProfileCount == 1 {
				fmt.Printf(" - 1 profile")
			} else {
				fmt.Printf(" - %d profiles", cli.ProfileCount)
			}
		}

		// Print connectivity status
		if cli.Connected {
			fmt.Printf(" ✓\n")
		} else {
			fmt.Printf(" ✗\n")
		}
	}
}

// printAWSDetails prints detailed AWS resource information
func printAWSDetails() {
	// Check if AWS CLI is available
	if !exists("aws") {
		fmt.Println("AWS CLI not found")
		return
	}

	// Get AWS version
	version := getCloudCLIVersion("aws")
	if version != "" {
		fmt.Printf("AWS CLI: %s\n\n", version)
	}

	// Get profiles
	allProfiles := getAWSProfiles()
	if len(allProfiles) == 0 {
		fmt.Println("No AWS profiles configured")
		return
	}

	// Filter profiles if --profile flag is specified
	var profiles []string
	if profileFlag != "" {
		// Check if specified profile exists
		profileExists := false
		for _, p := range allProfiles {
			if p == profileFlag {
				profileExists = true
				break
			}
		}
		if !profileExists {
			fmt.Printf("Profile '%s' not found. Available profiles: %v\n", profileFlag, allProfiles)
			return
		}
		profiles = []string{profileFlag}
	} else {
		profiles = allProfiles
	}

	fmt.Printf("Profiles: %d\n", len(profiles))

	// Get default profile if showing multiple profiles
	var defaultProfile string
	if len(profiles) > 1 {
		defaultProfile = getDefaultAWSProfile()
	}

	// Process each profile in parallel
	var wg sync.WaitGroup
	for _, profile := range profiles {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			isDefault := len(profiles) > 1 && p == defaultProfile
			printAWSProfileResources(p, isDefault)
		}(profile)
	}
	wg.Wait()
}

// printAWSProfileResources prints AWS resources for a specific profile
func printAWSProfileResources(profile string, isDefaultProfile bool) {
	// Check connectivity first
	connected := checkAWSProfileConnectivity(profile)
	connectStatus := "✓"
	if !connected {
		connectStatus = "✗"
	}

	// Build profile header
	profileHeader := fmt.Sprintf("Profile: %s", profile)
	if isDefaultProfile {
		profileHeader += " (default)"
	}
	profileHeader += fmt.Sprintf(" [%s]", connectStatus)
	fmt.Printf("\n%s\n", profileHeader)

	if !connected {
		fmt.Printf("  Unable to connect to AWS with this profile\n")
		return
	}

	// Get default region for this profile
	defaultRegion := getDefaultRegionForProfile(profile)
	if defaultRegion != "" {
		fmt.Printf("  Default region: %s\n", defaultRegion)
	}

	// Load AWS config for this profile
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		fmt.Printf("  Error loading config: %v\n", err)
		return
	}

	// Get regions to check
	regions := getAWSRegions(cfg)
	if regionFlag != "" {
		regions = []string{regionFlag}
	}

	// Query each region in parallel
	var wg sync.WaitGroup
	for _, region := range regions {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			isDefaultRegion := r == defaultRegion
			queryAWSRegion(profile, r, isDefaultRegion)
		}(region)
	}
	wg.Wait()
}

// getAWSRegions returns list of AWS regions to check
func getAWSRegions(cfg aws.Config) []string {
	// Try to get all regions dynamically from EC2
	ec2Client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	}

	result, err := ec2Client.DescribeRegions(context.TODO(), input)
	if err == nil && len(result.Regions) > 0 {
		regions := make([]string, 0, len(result.Regions))
		for _, region := range result.Regions {
			if region.RegionName != nil {
				regions = append(regions, *region.RegionName)
			}
		}
		return regions
	}

	// Fallback to comprehensive list of all known AWS regions
	return []string{
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
		"af-south-1",
		"ap-east-1",
		"ap-south-1",
		"ap-south-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ap-southeast-4",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ca-central-1",
		"eu-central-1",
		"eu-central-2",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-south-1",
		"eu-south-2",
		"eu-north-1",
		"il-central-1",
		"me-south-1",
		"me-central-1",
		"sa-east-1",
	}
}

// queryAWSRegion queries AWS Resource Groups Tagging API for resource counts in a region
func queryAWSRegion(profile, region string, isDefaultRegion bool) {
	// Load config with specific region
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return
	}

	// Create Resource Groups Tagging API client
	client := resourcegroupstaggingapi.NewFromConfig(cfg)

	// Get all resources in the region
	resourceCounts := make(map[string]int64)
	var paginationToken *string

	for {
		input := &resourcegroupstaggingapi.GetResourcesInput{
			PaginationToken:     paginationToken,
			ResourcesPerPage:    aws.Int32(100),
			ResourceTypeFilters: []string{}, // Empty means all resource types
		}

		resp, err := client.GetResources(context.TODO(), input)
		if err != nil {
			// Silently skip regions with errors (e.g., permission issues, service not available)
			return
		}

		// Count resources by type
		for _, resource := range resp.ResourceTagMappingList {
			if resource.ResourceARN != nil {
				resourceType := extractResourceTypeFromARN(*resource.ResourceARN)
				resourceCounts[resourceType]++
			}
		}

		// Check if there are more pages
		if resp.PaginationToken == nil || *resp.PaginationToken == "" {
			break
		}
		paginationToken = resp.PaginationToken
	}

	// Calculate total resources
	total := int64(0)
	for _, count := range resourceCounts {
		total += count
	}

	// Only print if there are resources
	if total > 0 {
		regionLabel := fmt.Sprintf("  Region: %s", region)
		if isDefaultRegion {
			regionLabel += " (default)"
		}
		regionLabel += fmt.Sprintf(" (Total: %d resources)", total)
		fmt.Printf("%s\n", regionLabel)

		for resourceType, count := range resourceCounts {
			if count > 0 {
				fmt.Printf("    %s: %d\n", resourceType, count)
			}
		}
	}
}

// extractResourceTypeFromARN extracts the resource type from an ARN
// ARN format: arn:partition:service:region:account-id:resource-type/resource-id
func extractResourceTypeFromARN(arn string) string {
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return "Unknown"
	}

	service := parts[2]
	resourcePart := parts[5]

	// Handle different ARN formats
	if strings.Contains(resourcePart, "/") {
		resourceTypeParts := strings.SplitN(resourcePart, "/", 2)
		return service + "::" + resourceTypeParts[0]
	}

	return service + "::" + resourcePart
}

// AWSResourceCount holds resource count information
type AWSResourceCount struct {
	ResourceType string
	Count        int64
}

// printCloudNativeForStatus prints cloud-native summary in status command format
func printCloudNativeForStatus() {
	clis := detectCloudCLIs()

	if len(clis) == 0 {
		// No output if no CLIs detected
		return
	}

	fmt.Println("Cloud Native:")
	for _, cli := range clis {
		// Print CLI name and version
		if cli.Version != "" {
			fmt.Printf("  %s (%s)", cli.Name, cli.Version)
		} else {
			fmt.Printf("  %s", cli.Name)
		}

		// For kubectl, also show kustomize version
		if cli.Name == "kubectl" && cli.KustomizeVersion != "" {
			fmt.Printf(" [kustomize: %s]", cli.KustomizeVersion)
		}

		// Print profile/context count (kubectl uses contexts, others use profiles)
		if cli.Name == "kubectl" {
			if cli.ProfileCount == 0 {
				fmt.Printf(" - 0 contexts")
			} else if cli.ProfileCount == 1 {
				fmt.Printf(" - 1 context")
			} else {
				fmt.Printf(" - %d contexts", cli.ProfileCount)
			}
		} else {
			if cli.ProfileCount == 0 {
				fmt.Printf(" - 0 profiles")
			} else if cli.ProfileCount == 1 {
				fmt.Printf(" - 1 profile")
			} else {
				fmt.Printf(" - %d profiles", cli.ProfileCount)
			}
		}

		// Print connectivity status
		if cli.Connected {
			fmt.Printf(" ✓\n")
		} else {
			fmt.Printf(" ✗\n")
		}
	}
	fmt.Println()
}
