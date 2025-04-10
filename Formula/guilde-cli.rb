class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.43.1"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.1/guilde-cli-darwin-arm64.zip"
      sha256 "b26e52141a8a483a5184bc60ff6276a766271fe1b640eeb5e82dcb2de8d98d29"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.1/guilde-cli-darwin-amd64.zip"
      sha256 "e93c185febe16d46e51df4e99be2267dfef007ef88451cb68549aff68eadf1e2"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.1/guilde-cli-linux-arm64.zip"
      sha256 "f119394ffb5aaa137721377459f88dcfd5c703cfee4d8a37a47c4fd5478a6247"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.1/guilde-cli-linux-amd64.zip"
      sha256 "b84c0144fee1df13211a863aff4091968f968b4a39f1f3c01ad97a70b622d341"
    end
  end

  def install
    bin.install "guilde"
  end
end

