class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.8"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.8/guilde-cli-darwin-arm64.zip"
      sha256 "20b43549e809ad7a6fb9ec23e3d2fe9110851f77abfb60f8a41dd5d5e737ebe7"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.8/guilde-cli-darwin-amd64.zip"
      sha256 "2afb10010e213b6b6b5fec3b5cf5d19d62b0843f77402c6a55b646c15cd87c8c"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.8/guilde-cli-linux-arm64.zip"
      sha256 "a4318681f61ec6b6dadc79a043d89f624aa4407f71eed61908b7651ea4e45d96"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.8/guilde-cli-linux-amd64.zip"
      sha256 "2500a12c169711fbc5fb183fb4a7207b4284ed376f667d104cc7669332fcc922"
    end
  end

  def install
    bin.install "guilde"
  end
end

