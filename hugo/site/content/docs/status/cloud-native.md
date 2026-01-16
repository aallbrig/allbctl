---
weight: 5
title: "Cloud Native"
---

# Status Cloud Native

Show cloud CLI tools and cloud resource information.

## Usage

```bash
# Show summary of all cloud CLIs
allbctl status cloud-native
allbctl status cn  # Short alias

# Show detailed AWS resources across all regions for all profiles
allbctl status cloud-native aws
allbctl status cn aws

# Show AWS resources in specific region for all profiles
allbctl status cn aws --region us-east-1

# Show AWS resources for specific profile across all regions
allbctl status cn aws --profile production

# Show AWS resources for specific profile and region
allbctl status cn aws --profile prod --region us-west-2
```

## Summary Output

The summary command shows:
- Detected cloud CLI tools (AWS CLI, gcloud, Azure CLI, kubectl)
- CLI versions (kubectl also shows kustomize version)
- Number of configured profiles/accounts (kubectl shows contexts instead of profiles)
- Connectivity status (✓ connected, ✗ not connected)

Example:
```
Cloud Native:
  kubectl (1.34.2) [kustomize: 5.7.1] - 0 contexts ✗
  aws (2.33.1) - 1 profile ✓
  gcloud (478.0.0) - 2 profiles ✓
```

### Connectivity Status

Each cloud CLI shows its connectivity status:
- **✓ (checkmark)**: CLI is connected and can communicate with its respective cloud service
  - **AWS**: At least one profile can authenticate (uses `aws sts get-caller-identity`)
  - **gcloud**: Has an active authenticated account
  - **Azure**: Can access account information
  - **kubectl**: Can connect to a Kubernetes cluster
- **✗ (x mark)**: CLI is not connected
  - **AWS**: No profiles can authenticate (all profiles fail connectivity check)
  - **kubectl**: No cluster configured or cannot reach cluster
  
In the AWS detailed view, each profile shows its individual connectivity status.

## AWS Detailed Output

The AWS detailed command queries AWS Resource Groups Tagging API to discover all resources across all services.

### Key Features

- **Universal discovery**: Uses Resource Groups Tagging API to find ALL AWS resources
- **No setup required**: Works without AWS Config service enabled
- **Comprehensive**: Detects resources across all AWS services (EC2, S3, Lambda, GameLift, RDS, etc.)
- **Auto-discovery**: Automatically checks all AWS regions
- **Smart filtering**: Only displays regions that contain resources (>=1)
- **Resource filtering**: Only displays resource types with count >= 1
- **Multi-profile support**: Processes all configured AWS profiles or filter to specific profile with `--profile`
- **Region filtering**: Query all regions or filter to specific region with `--region`
- **Parallel execution**: Queries regions concurrently for faster results

### Example Output

```
AWS CLI: 2.33.1

Profiles: 2

Profile: default (default) [✓]
  Default region: us-east-2
  Region: us-east-2 (default) (Total: 15 resources)
    gamelift::fleet: 3
    acm::certificate: 1
    lightsail::KeyPair: 2
    ssm::session: 4
    ecs::cluster: 1
  Region: us-east-1 (Total: 8 resources)
    ec2::instance: 5
    rds::db: 2

Profile: production [✓]
  Default region: us-west-2
  Region: us-west-2 (default) (Total: 45 resources)
    ec2::instance: 10
    rds::db: 5
    s3::production-bucket: 30

Profile: staging [✗]
  Unable to connect to AWS with this profile
```

**Note**: 
- **Profile indicators**: When multiple profiles exist, the default profile is marked with "(default)"
- **Region indicators**: Each profile's default region (from AWS config) is shown and marked with "(default)" in the resource list
- **Connectivity status**: 
  - `[✓]` - Profile is connected and can authenticate with AWS
  - `[✗]` - Profile cannot connect (credentials may be invalid or expired)

### Requirements

- AWS CLI must be installed and configured
- Appropriate IAM permissions to call `tag:GetResources`
- Resources are discovered via Resource Groups Tagging API (works for all taggable AWS resources)

**Note**: The Resource Groups Tagging API discovers resources that are taggable. Most AWS resources are included, but:
- Resources must be in a taggable service (most AWS services support tagging)
- Some very old resources or resources in certain states may not be returned
- To ensure all resources are discovered, consider adding tags to your resources

For a complete inventory, you can combine this with service-specific commands (e.g., `aws gamelift list-fleets`).

### Filtering Options

#### Region Filtering

When no `--region` flag is provided:
- All AWS regions are checked automatically
- Only regions with resources (total >= 1) are displayed
- Empty regions are silently skipped

When `--region` flag is provided:
- Only the specified region is queried
- If the region has no resources, no output is shown for that region

#### Profile Filtering

When no `--profile` flag is provided:
- All configured AWS profiles are queried
- Each profile's resources are listed separately

When `--profile` flag is provided:
- Only the specified profile is queried
- If the profile doesn't exist, an error message shows available profiles

#### Combined Filtering

You can combine both flags for precise queries:
```bash
allbctl status cn aws --profile production --region us-east-1
```
This queries only the production profile in the us-east-1 region.

## Detected CLIs

### AWS CLI (aws)
- Command: `aws`
- Profiles from: `~/.aws/config` and `~/.aws/credentials`
- Detailed view: Resource counts via AWS Resource Groups Tagging API

### Google Cloud CLI (gcloud)
- Command: `gcloud`
- Profiles from: `gcloud auth list`
- Detailed view: Implementation pending

### Azure CLI (az)
- Command: `az`
- Profiles from: `az account list`
- Detailed view: Implementation pending

### Kubernetes (kubectl)
- Command: `kubectl`
- Contexts from: `~/.kube/config`
- Shows: Client version, kustomize version, number of contexts
- Connectivity: Checks if kubectl can connect to a cluster
- Detailed view: Implementation pending

## Integration

The cloud-native summary is included in the main `allbctl status` command output.
