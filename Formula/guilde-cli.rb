class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.9"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.9/guilde-cli-darwin-arm64.zip"
      sha256 "ca788578620d3192298490eb5a8cf0176a278d9f09ee8596b34c06733e8070cf"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.9/guilde-cli-darwin-amd64.zip"
      sha256 "de4a8893d2068e3902f73d2f39e3ab416a04cad809b0bda32873fcd9e7b66aca"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.9/guilde-cli-linux-arm64.zip"
      sha256 "a7bc3a54a5787b7746d8c2518fd252b552c9d6e300085cdc9afe5a30a1086254"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.9/guilde-cli-linux-amd64.zip"
      sha256 "10d1fac0db3e3daeed420366f971b52219e1669b991dadc62d5a28e1f67197a4"
    end
  end

  def install
    bin.install "guilde"
  end
end

