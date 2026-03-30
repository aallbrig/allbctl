#!/usr/bin/env python3
"""Generate a .deb package for allbctl.

Usage:
    python3 generate-deb-package.py <version> <arch> <binary_path> <output_dir>

Creates:
    <output_dir>/allbctl_<version>_<arch>.deb

Example:
    python3 generate-deb-package.py 0.0.37 amd64 ./allbctl_linux_amd64 ./dist
"""
import os
import shutil
import stat
import subprocess
import sys

if len(sys.argv) != 5:
    print(f"Usage: {sys.argv[0]} <version> <arch> <binary_path> <output_dir>")
    sys.exit(1)

version, arch, binary_path, output_dir = sys.argv[1:]

pkg_name = "allbctl"
maintainer = "Andrew Allbright <andrew.allbright@gmail.com>"
description = (
    "CLI tool for managing and inspecting your development environment.\n"
    " allbctl (allbrightctl) provides system status information, runtime\n"
    " detection, project discovery, network diagnostics, and more."
)
homepage = "https://github.com/aallbrig/allbctl"

# Build directory structure
staging = os.path.join(output_dir, f"{pkg_name}_{version}_{arch}")
debian_dir = os.path.join(staging, "DEBIAN")
bin_dir = os.path.join(staging, "usr", "bin")
os.makedirs(debian_dir, exist_ok=True)
os.makedirs(bin_dir, exist_ok=True)

# Copy binary
dest_binary = os.path.join(bin_dir, pkg_name)
shutil.copy2(binary_path, dest_binary)
os.chmod(dest_binary, stat.S_IRWXU | stat.S_IRGRP | stat.S_IXGRP | stat.S_IROTH | stat.S_IXOTH)

# Calculate installed size in KB
installed_size = os.path.getsize(dest_binary) // 1024

# Write control file
control = f"""Package: {pkg_name}
Version: {version}
Architecture: {arch}
Maintainer: {maintainer}
Installed-Size: {installed_size}
Section: utils
Priority: optional
Homepage: {homepage}
Description: {description}
"""

with open(os.path.join(debian_dir, "control"), "w") as f:
    f.write(control)

# Build the .deb
os.makedirs(output_dir, exist_ok=True)
deb_filename = f"{pkg_name}_{version}_{arch}.deb"
deb_path = os.path.join(output_dir, deb_filename)
subprocess.run(["dpkg-deb", "--build", "--root-owner-group", staging, deb_path], check=True)

# Clean up staging directory
shutil.rmtree(staging)

print(f"Generated {deb_path}")
