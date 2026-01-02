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

  # Windows 10 Development Environment for allbctl testing
  config.vm.define "windows10" do |win|
    # Use Windows 10 box
    win.vm.box = "gusztavvargadr/windows-10"
    win.vm.hostname = "allbctl-win10-test"
    
    # VM provider settings
    win.vm.provider "virtualbox" do |vb|
      vb.name = "allbctl-windows10-test"
      vb.gui = true
      vb.memory = "4096"
      vb.cpus = 2
      # Enable clipboard sharing
      vb.customize ["modifyvm", :id, "--clipboard-mode", "bidirectional"]
      vb.customize ["modifyvm", :id, "--draganddrop", "bidirectional"]
    end
    
    # Network configuration
    win.vm.network "private_network", type: "dhcp"
    
    # Sync the built binary to the VM
    # Build with: make build-windows before running vagrant
    win.vm.synced_folder ".", "/vagrant", disabled: false
    
    # Provision script to setup test environment
    win.vm.provision "shell", privileged: false, inline: <<-SHELL
      Write-Host "Setting up Windows test environment for allbctl..."
      
      # Create a test directory
      $testDir = "C:\\allbctl-test"
      if (-not (Test-Path $testDir)) {
        New-Item -ItemType Directory -Path $testDir | Out-Null
        Write-Host "Created test directory: $testDir"
      }
      
      # Copy allbctl binary to test directory
      $binarySource = "C:\\vagrant\\allbctl_windows_amd64.exe"
      if (Test-Path $binarySource) {
        Copy-Item $binarySource "$testDir\\allbctl.exe" -Force
        Write-Host "Copied allbctl binary to $testDir"
      } else {
        Write-Host "Warning: allbctl binary not found at $binarySource"
        Write-Host "Build with 'make build-windows' before running vagrant up"
      }
      
      # Add test directory to PATH for current session
      $env:Path += ";$testDir"
      
      Write-Host ""
      Write-Host "========================================"
      Write-Host "Windows Test Environment Ready!"
      Write-Host "========================================"
      Write-Host ""
      Write-Host "To test allbctl:"
      Write-Host "1. Open PowerShell"
      Write-Host "2. cd C:\\allbctl-test"
      Write-Host "3. .\\allbctl.exe bootstrap status"
      Write-Host "4. .\\allbctl.exe bootstrap install"
      Write-Host "5. .\\allbctl.exe bootstrap status"
      Write-Host ""
      Write-Host "The binary is at: C:\\allbctl-test\\allbctl.exe"
      Write-Host "Source code is at: C:\\vagrant"
      Write-Host ""
    SHELL
  end
end
