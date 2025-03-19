class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.41.8"

  on_macos do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-darwin-arm64.zip"
      sha256 "c73652de2f7258547b18ea73c31af6414c0af96a29d35438876ec3df099a0512"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-darwin-amd64.zip"
      sha256 "0bab1f67671e195fd3ed8bd858f4ae720804115d4c267e72226022d3a561593d"
    end
  end

  on_linux do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-linux-arm64.zip"
      sha256 "cba01cd5ac9d10b10119c272a6728b519d86d7aadad97205ef605f0a52e68f26"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-linux-amd64.zip"
      sha256 "bbb401b9966da299d4cae51d64533af8f9b4b2e6f01f2fc09af82b123797e311"
    end
  end

  def install
    bin.install "guilde"
  end
end
