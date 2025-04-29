class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.0"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.0/guilde-cli-darwin-arm64.zip"
      sha256 "0b8cc2a7a750d44df7cea600989481af59ff2b78348be529b29b831cb3912610"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.0/guilde-cli-darwin-amd64.zip"
      sha256 "06ac40f9aa38e26ecfd10c5fcec9b72020bd010e3ba0e4d530e054076815fb28"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.0/guilde-cli-linux-arm64.zip"
      sha256 "a0873f51befc72def8e628eaf308c1e7255367d1b4b8a634cd884c619395681f"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.0/guilde-cli-linux-amd64.zip"
      sha256 "e767780386db875eac7badd4c82992aac28939a9e11ac03bbc6d9a2bf4284a1e"
    end
  end

  def install
    bin.install "guilde"
  end
end

