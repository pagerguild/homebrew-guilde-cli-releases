class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.5"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.5/guilde-cli-darwin-arm64.zip"
      sha256 "bf1f20d489561dd0afabfc81ec00a772a92e62c0f760fcb29cdfb905f93038e7"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.5/guilde-cli-darwin-amd64.zip"
      sha256 "81ba35c5ce8c82f56229175638fee944039462897a01305572d0e99cd1e1ed1e"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.5/guilde-cli-linux-arm64.zip"
      sha256 "20474c8d48f95ad059b2b6feb32d76731304236681c510fc0f6e352ab8a3c952"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.5/guilde-cli-linux-amd64.zip"
      sha256 "513c9a3b465a4ff9e89bb4119177787a9cdae684cc4c09f2aafbdf59be14a7ea"
    end
  end

  def install
    bin.install "guilde"
  end
end

