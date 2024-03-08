{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils = { url = "github:numtide/flake-utils"; };
  };
  outputs = { nixpkgs, flake-utils, ... }:
    (flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          name = "Shell with Go toolchain";
          packages = with pkgs; [ go gopls ];
        };
      })) // {
        nixosModules.default = import ./nix/nixos.nix;
        hmModules.default = import ./nix/home-manager.nix;
      };
}
