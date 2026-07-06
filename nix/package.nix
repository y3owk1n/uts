{
  fetchurl,
  gitUpdater,
  installShellFiles,
  stdenv,
  versionCheckHook,
  lib,
  buildGoModule,
  version ? "main",
  useZip ? false,
  commitHash ? null,
  writableTmpDirAsHomeHook,
  nix-update-script,
  unzip,
}:
if useZip then
  let
    # Determine architecture-specific details
    archInfo =
      {
        "aarch64-darwin" = {
          url = "https://github.com/y3owk1n/uts/releases/download/v${version}/uts-darwin-arm64.zip";
          # run `nix hash convert --hash-algo sha256 (nix-prefetch-url https://github.com/y3owk1n/uts/releases/download/v0.1.0/uts-darwin-arm64.zip)`
          sha256 = "";
        };
        "x86_64-darwin" = {
          url = "https://github.com/y3owk1n/uts/releases/download/v${version}/uts-darwin-amd64.zip";
          # run `nix hash convert --hash-algo sha256 (nix-prefetch-url https://github.com/y3owk1n/uts/releases/download/v0.1.0/uts-darwin-amd64.zip)`
          sha256 = "";
        };
        "aarch64-linux" = {
          url = "https://github.com/y3owk1n/uts/releases/download/v${version}/uts-linux-arm64.zip";
          # run `nix hash convert --hash-algo sha256 (nix-prefetch-url https://github.com/y3owk1n/uts/releases/download/v0.1.0/uts-linux-arm64.zip)`
          sha256 = "";
        };
        "x86_64-linux" = {
          url = "https://github.com/y3owk1n/uts/releases/download/v${version}/uts-linux-amd64.zip";
          # run `nix hash convert --hash-algo sha256 (nix-prefetch-url https://github.com/y3owk1n/uts/releases/download/v0.1.0/uts-linux-amd64.zip)`
          sha256 = "";
        };
      }
      .${stdenv.hostPlatform.system} or (throw "Unsupported system: ${stdenv.hostPlatform.system}");

  in
  stdenv.mkDerivation {
    pname = "uts";

    inherit version;

    src = fetchurl {
      url = archInfo.url;
      sha256 = archInfo.sha256;
    };

    unpackPhase = ''
      unzip $src
    '';

    nativeBuildInputs = [
      installShellFiles
      unzip
    ];

    installPhase = ''
      runHook preInstall
      mkdir -p $out/bin
      mv bin/uts $out/bin/uts
      mkdir -p $out/share/man/man1
      mv share/man/man1/*.1 $out/share/man/man1/
      runHook postInstall
    '';

    # only install completions on macOS
    # unable to make it work on Linux (do it manually please, sorry)
    postInstall = ''
      if ${
        lib.boolToString (
          stdenv.buildPlatform.canExecute stdenv.hostPlatform && stdenv.hostPlatform.isDarwin
        )
      }; then
        installShellCompletion --cmd uts \
              --bash <($out/bin/uts completion bash) \
              --fish <($out/bin/uts completion fish) \
              --zsh <($out/bin/uts completion zsh)
      fi
    '';

    doInstallCheck = true;
    nativeInstallCheckInputs = [
      versionCheckHook
    ];

    passthru.updateScript = gitUpdater {
      url = "https://github.com/y3owk1n/uts.git";
      rev-prefix = "v";
    };

    meta = with lib; {
      description = "One CLI for every format";
      homepage = "https://github.com/y3owk1n/uts";
      license = licenses.mit;
      platforms = platforms.darwin ++ platforms.linux;
      mainProgram = "uts";
    };
  }
else
  let
    shortHash = if commitHash != null then lib.substring 0 7 commitHash else null;

    pversion = "${version}${if shortHash != null then "-${shortHash}" else ""}";
  in
  # Build from source
  buildGoModule (finalAttrs: {
    pname = "uts";
    version = pversion;

    src = lib.cleanSource ../.;

    # run the following command to get the sha256 hash
    # `nix-shell -p go --run 'go mod vendor'`
    # `nix hash path vendor`
    # `rm -rf vendor`
    vendorHash = "sha256-OshIZkTb2wocyjZkCo7Ucwljxsedc1pGc0k/g3oOEa0=";

    ldflags = [
      "-s"
      "-w"
      "-X github.com/y3owk1n/uts/cmd.Version=${finalAttrs.version}"
    ]
    ++ lib.optionals (commitHash != null) [
      "-X github.com/y3owk1n/uts/cmd.GitCommit=${commitHash}"
    ];

    nativeBuildInputs = [
      installShellFiles
      writableTmpDirAsHomeHook
    ];

    subPackages = [ "." ];

    # Allow Go to use any available toolchain
    preBuild = ''
      export GOTOOLCHAIN=auto
    '';

    postInstall = ''
      # generate man pages
      mkdir -p $out/share/man/man1
      go run ./cmd/genman $out/share/man/man1

      # install shell completions
      if ${lib.boolToString (stdenv.buildPlatform.canExecute stdenv.hostPlatform)}; then
      	installShellCompletion --cmd neru \
      	--bash <($out/bin/neru completion bash) \
      	--fish <($out/bin/neru completion fish) \
      	--zsh <($out/bin/neru completion zsh)
      fi
    '';

    passthru = {
      updateScript = nix-update-script { };
    };

    meta = with lib; {
      description = "One CLI for every format";
      homepage = "https://github.com/y3owk1n/uts";
      license = licenses.mit;
      platforms = platforms.darwin ++ platforms.linux;
      mainProgram = "uts";
    };
  })
