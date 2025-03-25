class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.42.0"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.42.0/guilde-cli-darwin-arm64.zip"
      sha256 "1bb498995af9032fa1d45cfd6031dfa24285bb83e3a27b8f1da6d786328b3497"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.42.0/guilde-cli-darwin-amd64.zip"
      sha256 "bcdbcd195d2434633251e6eca621c3defeea0e57750d2fa99aa1582ac92e4f93"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.42.0/guilde-cli-linux-arm64.zip"
      sha256 "1bec677a12bfb1f2a99a1df47fc21ef3c1fde2ebead774e9cb3c69f229892660"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.42.0/guilde-cli-linux-amd64.zip"
      sha256 "1c1669e193605f4fc8a94b252571e47c654c8344771c547999fbd309abf14afc"
    end
  end

  def install
    bin.install "guilde"
  end
end

