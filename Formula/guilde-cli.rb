class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.4"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.4/guilde-cli-darwin-arm64.zip"
      sha256 "b8c4b5ecc16cb51e1b47f07ecb8c90f02ad3751af02424677ec44c30ed67fc88"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.4/guilde-cli-darwin-amd64.zip"
      sha256 "3892eed7411337a8b6cf0b94873ab40d882c25f9e5bb827da7a2aaa40390c442"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.4/guilde-cli-linux-arm64.zip"
      sha256 "93dbc327ff50ca3d15896e57e80c5efafbc912bb5a8e01744e433904d4a517ff"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.4/guilde-cli-linux-amd64.zip"
      sha256 "3816e74801c87bd94184f340bf569f797595120ea5850c1248fde25ba17d1eaa"
    end
  end

  def install
    bin.install "guilde"
  end
end

