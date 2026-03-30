#!/usr/bin/env python3
"""Update an APT repository with new .deb packages.

Usage:
    python3 update-apt-repo.py <repo_dir> <deb_file> [<deb_file> ...]

Updates the APT repository at <repo_dir> by:
  1. Copying .deb files into pool/main/a/allbctl/
  2. Regenerating Packages indices for each architecture
  3. Generating Release file
  4. Signing with GPG to produce InRelease and Release.gpg

Requires:
  - dpkg-scanpackages (from dpkg-dev)
  - apt-ftparchive (from apt-utils)
  - gpg with a signing key available (set GPG_KEY_ID env var)
"""
import gzip
import os
import shutil
import subprocess
import sys

if len(sys.argv) < 3:
    print(f"Usage: {sys.argv[0]} <repo_dir> <deb_file> [<deb_file> ...]")
    sys.exit(1)

repo_dir = sys.argv[1]
deb_files = sys.argv[2:]

gpg_key_id = os.environ.get("GPG_KEY_ID", "")
if not gpg_key_id:
    print("Error: GPG_KEY_ID environment variable must be set")
    sys.exit(1)

pool_dir = os.path.join(repo_dir, "pool", "main", "a", "allbctl")
dists_dir = os.path.join(repo_dir, "dists", "stable")
os.makedirs(pool_dir, exist_ok=True)

# Copy .deb files into pool
for deb in deb_files:
    dest = os.path.join(pool_dir, os.path.basename(deb))
    shutil.copy2(deb, dest)
    print(f"Copied {deb} -> {dest}")

# Determine architectures from the .deb filenames
architectures = set()
for deb in deb_files:
    basename = os.path.basename(deb)
    # Expected format: allbctl_<version>_<arch>.deb
    parts = basename.replace(".deb", "").split("_")
    if len(parts) >= 3:
        architectures.add(parts[-1])

if not architectures:
    architectures = {"amd64", "arm64"}

# Generate Packages index for each architecture
for arch in architectures:
    arch_dir = os.path.join(dists_dir, "main", f"binary-{arch}")
    os.makedirs(arch_dir, exist_ok=True)

    packages_path = os.path.join(arch_dir, "Packages")
    result = subprocess.run(
        ["dpkg-scanpackages", "--arch", arch, "pool/"],
        capture_output=True, text=True, cwd=repo_dir, check=True,
    )
    with open(packages_path, "w") as f:
        f.write(result.stdout)

    # Compress
    with open(packages_path, "rb") as f_in:
        with gzip.open(packages_path + ".gz", "wb") as f_out:
            f_out.write(f_in.read())

    print(f"Generated Packages index for {arch}")

# Generate Release file
arch_list = " ".join(sorted(architectures))
release_conf = f"""APT::FTPArchive::Release::Origin "aallbrig";
APT::FTPArchive::Release::Label "allbctl";
APT::FTPArchive::Release::Suite "stable";
APT::FTPArchive::Release::Codename "stable";
APT::FTPArchive::Release::Architectures "{arch_list}";
APT::FTPArchive::Release::Components "main";
"""
conf_path = os.path.join(repo_dir, "apt-ftparchive-release.conf")
with open(conf_path, "w") as f:
    f.write(release_conf)

release_path = os.path.join(dists_dir, "Release")
result = subprocess.run(
    ["apt-ftparchive", "-c", conf_path, "release", dists_dir],
    capture_output=True, text=True, check=True,
)
with open(release_path, "w") as f:
    f.write(result.stdout)
os.remove(conf_path)
print("Generated Release file")

# Sign: InRelease (inline signed)
inrelease_path = os.path.join(dists_dir, "InRelease")
subprocess.run(
    [
        "gpg", "--default-key", gpg_key_id,
        "--batch", "--yes", "--armor",
        "--clearsign",
        "--output", inrelease_path,
        release_path,
    ],
    check=True,
)
print("Generated InRelease (inline-signed)")

# Sign: Release.gpg (detached signature)
release_gpg_path = os.path.join(dists_dir, "Release.gpg")
subprocess.run(
    [
        "gpg", "--default-key", gpg_key_id,
        "--batch", "--yes", "--armor",
        "--detach-sign",
        "--output", release_gpg_path,
        release_path,
    ],
    check=True,
)
print("Generated Release.gpg (detached signature)")

print("APT repository updated successfully")
