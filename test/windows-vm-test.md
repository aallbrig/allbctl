# Windows VM Testing Guide

This guide explains how to test allbctl on a Windows 10 VM using Vagrant.

## Prerequisites

1. Install [VirtualBox](https://www.virtualbox.org/wiki/Downloads)
2. Install [Vagrant](https://www.vagrantup.com/downloads)

## Setup

1. Build the Windows binary:
   ```bash
   make build-windows
   ```

2. Start the Windows VM:
   ```bash
   vagrant up windows10
   ```

   This will:
   - Download Windows 10 box (first time only, ~10GB)
   - Create a VM with 4GB RAM and 2 CPUs
   - Copy the allbctl binary to `C:\allbctl-test\`
   - Setup the test environment

3. Wait for the VM to boot (GUI window will appear)

## Testing

Once the VM is running:

1. Log into Windows (credentials from the box, usually vagrant/vagrant)
2. Open PowerShell
3. Navigate to the test directory:
   ```powershell
   cd C:\allbctl-test
   ```

4. Run the bootstrap test sequence:
   ```powershell
   # Check initial status (should show missing items)
   .\allbctl.exe bootstrap status
   
   # Run installation
   .\allbctl.exe bootstrap install
   
   # Check status again (should show installed items)
   .\allbctl.exe bootstrap status
   ```

## Expected Test Results

### Initial `bootstrap status` should show:
- ❌ `NOT FOUND` for git
- ❌ `NOT FOUND` for gh
- ❌ `SSH KEY NOT FOUND` for SSH key
- ❌ `NOT CLONED` for dotfiles
- ❌ Directory `C:\Users\vagrant\src` does NOT exist

### After `bootstrap install` should show:
- ✅ git installed
- ✅ gh installed
- ✅ SSH key generated
- ✅ Directory `C:\Users\vagrant\src` created
- ✅ Dotfiles cloned (if git/gh work)

### Final `bootstrap status` should show:
- ✅ `INSTALLED` for git
- ✅ `INSTALLED` for gh
- ✅ `SSH KEY FOUND` for SSH key
- ✅ `CLONED` for dotfiles
- ✅ `PRESENT` for directory

## Known Issues to Test

1. **Directory detection**: Verify `~/src` directory status is accurate
2. **winget prompts**: Check if winget requires user input and handle appropriately
3. **Path variables**: Ensure `%USERPROFILE%` is used correctly for home directory

## Cleanup

```bash
# Shutdown VM
vagrant halt windows10

# Destroy VM completely
vagrant destroy windows10 -f
```

## VM Management

```bash
# SSH into VM (limited on Windows, use GUI)
vagrant rdp windows10

# Reload VM with new provisioning
vagrant reload windows10 --provision

# Check VM status
vagrant status
```

## Rebuilding After Code Changes

1. Rebuild Windows binary:
   ```bash
   make build-windows
   ```

2. Reload the VM:
   ```bash
   vagrant reload windows10 --provision
   ```

3. The new binary will be copied to `C:\allbctl-test\allbctl.exe`

## Troubleshooting

- **Box download slow**: The Windows 10 box is large (~10GB), first download takes time
- **VM won't start**: Check VirtualBox is installed and virtualization is enabled in BIOS
- **Binary not found**: Ensure you ran `make build-windows` before `vagrant up`
- **Winget issues**: winget may require Windows updates or manual acceptance of terms
