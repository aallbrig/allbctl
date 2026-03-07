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
    win.vm.boot_timeout = 600  # Windows is slow to boot; default 300s is not enough

    win.vm.provider "virtualbox" do |vb|
      vb.name = "allbctl-windows10-test"
      vb.gui = false   # headless; set VAGRANT_GUI=1 to enable desktop
      vb.memory = "4096"
      vb.cpus = 2
    end

    # Sync the project root into the VM so the binary is available at C:\vagrant
    win.vm.synced_folder ".", "/vagrant", disabled: false

    # Setup: copy binary + add to machine PATH, create a test git repo
    # for testing `allbctl status projects`
    win.vm.provision "shell", privileged: true,
      path: "scripts/vagrant-windows-setup.ps1"

    # Smoke-test provision: runs allbctl commands and prints results
    win.vm.provision "shell", run: "never", name: "smoke-test", privileged: false,
      path: "scripts/vagrant-windows-smoke-test.ps1"
  end
end
