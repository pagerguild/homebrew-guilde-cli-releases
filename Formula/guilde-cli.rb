class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.3"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.3/guilde-cli-darwin-arm64.zip"
      sha256 "c1becb07e65f13a731d8c3f5b94948b115df447fac9b4006c537ef3eac6b30fe"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.3/guilde-cli-darwin-amd64.zip"
      sha256 "7acbec0bc3b572477da0bff2b318ef80ba4fec0f30216825b3e0a01bc58afb66"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.3/guilde-cli-linux-arm64.zip"
      sha256 "faf4185dd86e2014af89e83c8838798c7e97e0a6c860d5e15c2c28370e300975"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.3/guilde-cli-linux-amd64.zip"
      sha256 "18d5b021f827d7727acd125cabab176cc298eff03daa4cb76f1909e709cc6751"
    end
  end

  def install
    bin.install "guilde"
  end
end

