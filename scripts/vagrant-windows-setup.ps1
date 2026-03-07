$testDir = "C:\allbctl-test"
if (-not (Test-Path $testDir)) { New-Item -ItemType Directory -Path $testDir | Out-Null }

$binarySource = "C:\vagrant\allbctl_windows_amd64.exe"
if (Test-Path $binarySource) {
    Copy-Item $binarySource "$testDir\allbctl.exe" -Force
    Write-Host "OK: copied allbctl.exe to $testDir"
} else {
    Write-Error "MISSING: $binarySource - run 'make build-windows' on the host first"
    exit 1
}

# Add to machine-wide PATH so all future sessions have it
$machinePath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
if ($machinePath -notlike "*$testDir*") {
    [System.Environment]::SetEnvironmentVariable("Path", "$machinePath;$testDir", "Machine")
    Write-Host "OK: added $testDir to machine PATH"
}

# Create a test git repo so status projects has something to scan
$srcDir = "C:\Users\vagrant\src\allbctl"
if (-not (Test-Path $srcDir)) {
    New-Item -ItemType Directory -Path $srcDir | Out-Null
    & git init $srcDir
    Write-Host "OK: created test git repo at $srcDir"
}
