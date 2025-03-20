class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.41.12"

  on_macos do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.12-darwin-arm64.zip"
      sha256 "44449cbf714b8e68f89031b2cd25b2b1ce7d4cce03352c270e4aac3335cbc8c2"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.12-darwin-amd64.zip"
      sha256 "07bb59cf12e12b23b5b6e7565cbba5971fd75471625b8f16e628d06ba8b7db31"
    end
  end

  on_linux do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.12-linux-arm64.zip"
      sha256 "fbe86dd541a343ac7ecbbf1522df4ec66478c39f6d9dae4f25b395ee41a1a957"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.12-linux-amd64.zip"
      sha256 "fb9a79a7d8adc10f1efc2402e22e2d98eaf1601d966b0971cf4bfee03c670e75"
    end
  end

  def install
    bin.install "guilde"
  end
end
