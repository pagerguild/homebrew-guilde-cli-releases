class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.43.2"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.2/guilde-cli-darwin-arm64.zip"
      sha256 "5a93386a84a236e00b21bb953db14386c8761980508a04c541eb8e5d7888fe2d"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.2/guilde-cli-darwin-amd64.zip"
      sha256 "cb85d865677c1075947bf50bc5b58f183dd860f92217a34ba18ad422be7994d2"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.2/guilde-cli-linux-arm64.zip"
      sha256 "5369e05a547c4214121dc899cd8c4e0555b66358b66854ffef841e38b3903e90"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.2/guilde-cli-linux-amd64.zip"
      sha256 "bc2f44af60f95895cccad0bf024e70699e655ad261780fafbd6c6b5ecf969acf"
    end
  end

  def install
    bin.install "guilde"
  end
end

