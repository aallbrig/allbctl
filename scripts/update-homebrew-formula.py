#!/usr/bin/env python3
"""Update the allbctl Homebrew formula in aallbrig/homebrew-tap.

Usage:
    python3 update-homebrew-formula.py <version> <base_url> \
        <darwin_amd64_sha> <darwin_arm64_sha> \
        <linux_amd64_sha> <linux_arm64_sha>
"""
import sys

version, base_url, d_amd64, d_arm64, l_amd64, l_arm64 = sys.argv[1:]

formula = (
    'class Allbctl < Formula\n'
    '  desc "CLI tool for managing and inspecting your development environment"\n'
    '  homepage "https://aallbrig.github.io/allbctl"\n'
    f'  version "{version}"\n'
    '  license "MIT"\n'
    '\n'
    '  on_macos do\n'
    '    on_intel do\n'
    f'      url "{base_url}/allbctl-darwin-amd64.tar.gz"\n'
    f'      sha256 "{d_amd64}"\n'
    '\n'
    '      def install\n'
    '        bin.install "allbctl_darwin_amd64" => "allbctl"\n'
    '      end\n'
    '    end\n'
    '    on_arm do\n'
    f'      url "{base_url}/allbctl-darwin-arm64.tar.gz"\n'
    f'      sha256 "{d_arm64}"\n'
    '\n'
    '      def install\n'
    '        bin.install "allbctl_darwin_arm64" => "allbctl"\n'
    '      end\n'
    '    end\n'
    '  end\n'
    '\n'
    '  on_linux do\n'
    '    on_intel do\n'
    f'      url "{base_url}/allbctl-linux-amd64.tar.gz"\n'
    f'      sha256 "{l_amd64}"\n'
    '\n'
    '      def install\n'
    '        bin.install "allbctl_linux_amd64" => "allbctl"\n'
    '      end\n'
    '    end\n'
    '    on_arm do\n'
    f'      url "{base_url}/allbctl-linux-arm64.tar.gz"\n'
    f'      sha256 "{l_arm64}"\n'
    '\n'
    '      def install\n'
    '        bin.install "allbctl_linux_arm64" => "allbctl"\n'
    '      end\n'
    '    end\n'
    '  end\n'
    '\n'
    '  test do\n'
    '    assert_match "v#{version}", shell_output("#{bin}/allbctl version")\n'
    '  end\n'
    'end\n'
)

with open("homebrew-tap/Formula/allbctl.rb", "w") as f:
    f.write(formula)

print(f"Updated formula to v{version}")
