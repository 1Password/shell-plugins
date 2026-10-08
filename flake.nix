{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
  };
  outputs = { nixpkgs, systems, ... }:
    let
      forAllSystems = function:
        nixpkgs.lib.genAttrs (import systems) (
          system: function nixpkgs.legacyPackages.${system}
        );
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          name = "Shell with Go toolchain";
          packages = with pkgs; [ go gopls ];
        };
      });

      nixosModules.default = import ./nix/nixos.nix;
      hmModules.default = import ./nix/home-manager.nix;
    };
}
