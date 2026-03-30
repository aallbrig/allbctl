#!/usr/bin/env python3
"""Generate Chocolatey package files for allbctl.

Usage:
    python3 generate-chocolatey-package.py <version> <base_url> \
        <windows_amd64_sha> <windows_arm64_sha>

Creates:
    chocolatey/allbctl.nuspec
    chocolatey/tools/chocolateyInstall.ps1
    chocolatey/tools/chocolateyUninstall.ps1
"""
import os
import sys

version, base_url, w_amd64_sha, w_arm64_sha = sys.argv[1:]

out_dir = "chocolatey"
tools_dir = os.path.join(out_dir, "tools")
os.makedirs(tools_dir, exist_ok=True)

nuspec = f"""<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.chocolatey.org/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>allbctl</id>
    <version>{version}</version>
    <title>allbctl</title>
    <authors>Andrew Allbright</authors>
    <owners>aallbrig</owners>
    <projectUrl>https://github.com/aallbrig/allbctl</projectUrl>
    <docsUrl>https://aallbrig.github.io/allbctl/</docsUrl>
    <bugTrackerUrl>https://github.com/aallbrig/allbctl/issues</bugTrackerUrl>
    <projectSourceUrl>https://github.com/aallbrig/allbctl</projectSourceUrl>
    <packageSourceUrl>https://github.com/aallbrig/allbctl</packageSourceUrl>
    <licenseUrl>https://github.com/aallbrig/allbctl/blob/main/LICENSE</licenseUrl>
    <requireLicenseAcceptance>false</requireLicenseAcceptance>
    <tags>allbctl cli devtools developer-tools sysadmin</tags>
    <summary>CLI tool for managing and inspecting your development environment</summary>
    <description>allbctl (allbrightctl) is a command-line interface for inspecting and managing your development environment. It provides system status information, runtime detection, project discovery, network diagnostics, and more.</description>
  </metadata>
</package>
"""

install_ps1 = f"""$ErrorActionPreference = 'Stop'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{{
  packageName    = 'allbctl'
  unzipLocation  = $toolsDir
  url64bit       = '{base_url}/allbctl-windows-amd64.zip'
  checksum64     = '{w_amd64_sha}'
  checksumType64 = 'sha256'
}}
Install-ChocolateyZipPackage @packageArgs

# Rename the extracted binary so Chocolatey shims it as 'allbctl'
$exe = Join-Path $toolsDir 'allbctl_windows_amd64.exe'
$target = Join-Path $toolsDir 'allbctl.exe'
if (Test-Path $exe) {{
  Move-Item -Force $exe $target
}}

# Suppress shimming for the original name if it somehow persists
$ignore = Join-Path $toolsDir 'allbctl_windows_amd64.exe.ignore'
Set-Content -Path $ignore -Value ''
"""

uninstall_ps1 = """$ErrorActionPreference = 'Stop'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$exe = Join-Path $toolsDir 'allbctl.exe'
if (Test-Path $exe) {
  Remove-Item -Force $exe
}
"""

with open(os.path.join(out_dir, "allbctl.nuspec"), "w") as f:
    f.write(nuspec)
with open(os.path.join(tools_dir, "chocolateyInstall.ps1"), "w") as f:
    f.write(install_ps1)
with open(os.path.join(tools_dir, "chocolateyUninstall.ps1"), "w") as f:
    f.write(uninstall_ps1)

print(f"Generated Chocolatey package files for v{version}")
