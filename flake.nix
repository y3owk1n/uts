{
  description = "Single CLI for files compression and conversion";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      eachSystem = nixpkgs.lib.genAttrs [
        "aarch64-darwin"
        "x86_64-darwin"
        "aarch64-linux"
        "x86_64-linux"
      ];

      latestVersion = "0.1.3";

      makeUtsPackage =
        pkgs: version: useZip: commitHash:
        pkgs.callPackage ./nix/package.nix {
          inherit version useZip commitHash;
        };
    in
    {
      overlays.default = final: prev: {
        uts = makeUtsPackage final latestVersion true null;
        uts-source = makeUtsPackage final "main" false (self.rev or self.dirtyRev or "unknown");
      };
      packages = eachSystem (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [ self.overlays.default ];
          };
        in
        {
          default = makeUtsPackage pkgs latestVersion true null;
          source = makeUtsPackage pkgs "main" false (self.rev or self.dirtyRev or "unknown");
        }
      );

      homeManagerModules.default = import ./nix/home.nix;
    };
}
