# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  # Ubuntu 24.04 LTS Development Environment
  config.vm.define "ubuntu" do |ubuntu|
    ubuntu.vm.box = "ubuntu/jammy64"  # Ubuntu 22.04 LTS (noble64 not available yet)
    ubuntu.vm.hostname = "allbctl-ubuntu-test"
    
    ubuntu.vm.provider "virtualbox" do |vb|
      vb.name = "allbctl-ubuntu-test"
      vb.memory = "2048"
      vb.cpus = 2
    end
    
    ubuntu.vm.network "private_network", type: "dhcp"
    ubuntu.vm.synced_folder ".", "/vagrant"
    
    ubuntu.vm.provision "shell", inline: <<-SHELL
      echo "Setting up Ubuntu test environment for allbctl..."
      
      # Update package lists
      apt-get update -qq
      
      # Install basic dependencies
      apt-get install -y git curl
      
      # Copy allbctl binary
      if [ -f /vagrant/bin/allbctl ]; then
        cp /vagrant/bin/allbctl /usr/local/bin/
        chmod +x /usr/local/bin/allbctl
        echo "✓ Copied allbctl binary to /usr/local/bin/"
      else
        echo "⚠ Warning: allbctl binary not found at /vagrant/bin/allbctl"
        echo "  Build with 'make build' before running vagrant up"
      fi
      
      echo ""
      echo "========================================"
      echo "Ubuntu Test Environment Ready!"
      echo "========================================"
      echo ""
      echo "To test allbctl:"
      echo "  vagrant ssh ubuntu"
      echo "  allbctl bootstrap status"
      echo "  allbctl bootstrap install"
      echo "  allbctl bootstrap status"
      echo ""
    SHELL
  end

  # Arch Linux Development Environment
  config.vm.define "arch" do |arch|
    arch.vm.box = "archlinux/archlinux"
    arch.vm.hostname = "allbctl-arch-test"
    
    arch.vm.provider "virtualbox" do |vb|
      vb.name = "allbctl-arch-test"
      vb.memory = "2048"
      vb.cpus = 2
    end
    
    arch.vm.network "private_network", type: "dhcp"
    arch.vm.synced_folder ".", "/vagrant"
    
    arch.vm.provision "shell", inline: <<-SHELL
      echo "Setting up Arch Linux test environment for allbctl..."
      
      # Update system
      pacman -Syu --noconfirm --needed
      
      # Install basic dependencies
      pacman -S --noconfirm --needed git curl
      
      # Copy allbctl binary
      if [ -f /vagrant/bin/allbctl ]; then
        cp /vagrant/bin/allbctl /usr/local/bin/
        chmod +x /usr/local/bin/allbctl
        echo "✓ Copied allbctl binary to /usr/local/bin/"
      else
        echo "⚠ Warning: allbctl binary not found at /vagrant/bin/allbctl"
        echo "  Build with 'make build' before running vagrant up"
      fi
      
      echo ""
      echo "========================================"
      echo "Arch Linux Test Environment Ready!"
      echo "========================================"
      echo ""
      echo "To test allbctl:"
      echo "  vagrant ssh arch"
      echo "  allbctl bootstrap status"
      echo "  allbctl bootstrap install"
      echo "  allbctl bootstrap status"
      echo ""
    SHELL
  end

  # Windows 10 Development/Test Environment for allbctl
  # Prerequisites: run `make build-windows` first to produce allbctl_windows_amd64.exe
  # Boot:  vagrant up windows10
  # Test:  vagrant powershell windows10 -c "C:\\allbctl-test\\allbctl.exe version"
  # Shell: vagrant powershell windows10
  config.vm.define "windows10" do |win|
    win.vm.box = "gusztavvargadr/windows-10"
    win.vm.hostname = "allbctl-win10-test"

    win.vm.provider "virtualbox" do |vb|
      vb.name = "allbctl-windows10-test"
      vb.gui = false   # headless; set VAGRANT_GUI=1 to enable desktop
      vb.memory = "4096"
      vb.cpus = 2
    end

    # Sync the project root into the VM so the binary is available at C:\vagrant
    win.vm.synced_folder ".", "/vagrant", disabled: false

    # Setup: copy binary + add to machine PATH, create a fake ~/src/allbctl git repo
    # for testing `allbctl status projects`
    win.vm.provision "shell", privileged: true, inline: <<-SHELL
      $testDir = "C:\\allbctl-test"
      if (-not (Test-Path $testDir)) { New-Item -ItemType Directory -Path $testDir | Out-Null }

      $binarySource = "C:\\vagrant\\allbctl_windows_amd64.exe"
      if (Test-Path $binarySource) {
        Copy-Item $binarySource "$testDir\\allbctl.exe" -Force
        Write-Host "OK: copied allbctl.exe to $testDir"
      } else {
        Write-Error "MISSING: $binarySource — run 'make build-windows' on the host first"
        exit 1
      }

      # Add to machine-wide PATH so all future sessions have it
      $machinePath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
      if ($machinePath -notlike "*$testDir*") {
        [System.Environment]::SetEnvironmentVariable("Path", "$machinePath;$testDir", "Machine")
        Write-Host "OK: added $testDir to machine PATH"
      }

      # Create a fake src/allbctl git repo so status projects has something to scan
      $srcDir = "C:\\Users\\vagrant\\src\\allbctl"
      if (-not (Test-Path $srcDir)) {
        New-Item -ItemType Directory -Path $srcDir | Out-Null
        & git init $srcDir
        Write-Host "OK: created test git repo at $srcDir"
      }
    SHELL

    # Smoke-test provision: runs allbctl commands and prints results
    win.vm.provision "shell", run: "never", name: "smoke-test", privileged: false, inline: <<-SHELL
      $allbctl = "C:\\allbctl-test\\allbctl.exe"
      $pass = 0; $fail = 0

      function Run-Test($label, $args) {
        Write-Host ""
        Write-Host "=== $label ===" -ForegroundColor Cyan
        $output = & $allbctl @args 2>&1
        $code = $LASTEXITCODE
        Write-Host $output
        if ($code -ne 0) {
          Write-Host "FAIL (exit $code)" -ForegroundColor Red
          $script:fail++
        } else {
          Write-Host "PASS" -ForegroundColor Green
          $script:pass++
        }
      }

      Run-Test "version"          @("version")
      Run-Test "help"             @("--help")
      Run-Test "status (full)"    @("status")
      Run-Test "status projects"  @("status", "projects")
      Run-Test "status runtimes"  @("status", "runtimes")
      Run-Test "bootstrap status" @("bootstrap", "status")

      Write-Host ""
      Write-Host "========================================"
      Write-Host "Results: $pass passed, $fail failed"
      if ($fail -gt 0) { Write-Host "OVERALL: FAIL" -ForegroundColor Red; exit 1 }
      else             { Write-Host "OVERALL: PASS" -ForegroundColor Green }
    SHELL
  end
end
