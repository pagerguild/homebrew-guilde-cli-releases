class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.41.10"

  on_macos do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.10-darwin-arm64.zip"
      sha256 "10a31fd6bf245d30e2ad646f8443fc4770f7e62fc2b82646a46ba1debc5026ad"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.10-darwin-amd64.zip"
      sha256 "728bd7a19c1e1fc760646d738e7dff132555a24da0aa8e256fc15013b7c718fc"
    end
  end

  on_linux do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.10-linux-arm64.zip"
      sha256 "a57cd3dd9c9ea1817816dc4740bf128c365752a6bffae647fbd2bfa1e8801a1b"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.10-linux-amd64.zip"
      sha256 "1bc15218a82a912277c64bcbc1d5fe960f28b36993f21016969b3a92f0200edc"
    end
  end

  def install
    bin.install "guilde"
  end
end
