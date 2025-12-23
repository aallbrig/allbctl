# Windows Testing Environment - Summary

## Changes Made

This commit adds a Windows 10 VM testing environment using Vagrant to test allbctl bootstrap functionality on Windows.

### Files Created

1. **Vagrantfile** - Vagrant configuration for Windows 10 VM
   - Uses `gusztavvargadr/windows-10` box
   - 4GB RAM, 2 CPUs
   - GUI enabled for PowerShell testing
   - Auto-provisions test environment
   - Copies built binary to `C:\allbctl-test\allbctl.exe`

2. **test/windows-vm-test.md** - Comprehensive testing guide
   - Prerequisites (VirtualBox, Vagrant)
   - Setup instructions
   - Testing procedure
   - Expected results
   - Known issues to test
   - VM management commands

### Files Modified

1. **.gitignore** - Added Vagrant entries
   - `.vagrant/` directory
   - `*.box` files

2. **README.md** - Added Windows VM Testing section
   - Quick start commands
   - Reference to detailed guide

## Testing Workflow

```bash
# 1. Build Windows binary
make build-windows

# 2. Start Windows VM
vagrant up windows10

# 3. In VM PowerShell:
cd C:\allbctl-test
.\allbctl.exe bootstrap status    # Should show missing items
.\allbctl.exe bootstrap install   # Should install components
.\allbctl.exe bootstrap status    # Should show installed items

# 4. Cleanup
vagrant halt windows10      # Stop VM
vagrant destroy windows10   # Remove VM
```

## Known Issues to Address (from user's report)

Based on the Windows PowerShell output provided, these issues need to be fixed:

### 1. Directory Detection False Positive
**Issue**: `~/src` directory was reported as `PRESENT` when it didn't exist
```
Expected Directories
-----
PRESENT C:\Users\aallb\src
```
But `ls ~/src` showed: `Cannot find path 'C:\Users\aallb\src' because it does not exist.`

**Location**: `pkg/osagnostic/ExpectedDirectory.go` line 25-33
**Problem**: The `Validate()` method checks `os.Stat()` but may have logic error

### 2. Winget Interactive Prompts
**Issue**: winget waits for user input during installation
```
Do you agree to all the source agreements terms?
[Y] Yes  [N] No: An unexpected error occurred while executing the command:
0x8a150042 : Error reading input in prompt
```

**Location**: `pkg/osagnostic/InstallableCommand.go` line 150-160
**Fix Needed**: Add `--accept-source-agreements` flag to winget commands
```powershell
winget install Git.Git --accept-source-agreements --accept-package-agreements
```

### 3. Windows Environment Variable
**Question from user**: What's the correct Windows variable for home directory?
- `%USERPROFILE%` is correct (e.g., `C:\Users\aallb`)
- `~` should expand to this in PowerShell
- Go's `os.UserHomeDir()` handles this correctly

## Next Steps

1. **Fix directory detection** - Debug why `os.Stat()` returns false positive
2. **Add winget flags** - Update `InstallableCommand.installWithPackageManager()` to add acceptance flags for winget
3. **Test with Vagrant** - Use this VM to validate fixes work correctly
4. **Consider interactive mode** - Maybe allow user to respond to prompts OR add `--yes` flag to allbctl

## Why Vagrant?

- **Reproducible**: Same environment every time
- **Safe**: Isolated from host system
- **Fast**: Can destroy and recreate in minutes
- **Cross-platform**: Works on Linux, macOS, Windows hosts
- **Realistic**: Real Windows 10 environment, not WSL

## Box Information

The `gusztavvargadr/windows-10` box:
- ~10GB download (first time only)
- Windows 10 Pro
- Already includes common tools
- Pre-activated for development
- Regular updates from maintainer
