class GuildeCli < Formula
  desc "Guilde CLI tool"
  homepage "https://github.com/pagerguild/guilde-cli-releases"
  version "0.1.0" # This will need to be updated with the actual version
  
  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/pagerguild/guilde-cli-releases/releases/download/v#{version}/guilde-cli-darwin-arm64.tar.gz"
    sha256 "SHA256_TO_BE_UPDATED" # This will need to be updated with actual SHA
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/pagerguild/guilde-cli-releases/releases/download/v#{version}/guilde-cli-darwin-amd64.tar.gz"
    sha256 "SHA256_TO_BE_UPDATED" # This will need to be updated with actual SHA
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/pagerguild/guilde-cli-releases/releases/download/v#{version}/guilde-cli-linux-arm64.tar.gz"
    sha256 "SHA256_TO_BE_UPDATED" # This will need to be updated with actual SHA
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/pagerguild/guilde-cli-releases/releases/download/v#{version}/guilde-cli-linux-amd64.tar.gz"
    sha256 "SHA256_TO_BE_UPDATED" # This will need to be updated with actual SHA
  else
    odie "Unsupported architecture: #{Hardware::CPU.arch}"
  end
  
  def install
    bin.install "bin/guilde-cli"
  end

  test do
    system "#{bin}/guilde-cli", "--version"
  end
end
