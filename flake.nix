{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };
  outputs =
    {
      nixpkgs,
      systems,
      treefmt-nix,
      ...
    }:
    let
      forAllSystems =
        function: nixpkgs.lib.genAttrs (import systems) (system: function nixpkgs.legacyPackages.${system});

      treefmtEval = forAllSystems (pkgs: treefmt-nix.lib.evalModule pkgs ./treefmt.nix);
    in
    {
      formatter = forAllSystems (pkgs: treefmtEval.${pkgs.system}.config.build.wrapper);

      checks = forAllSystems (pkgs: {
        formatting = treefmtEval.${pkgs.system}.config.build.check (pkgs.path or ./.);
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          name = "Shell with Go toolchain";
          packages = with pkgs; [
            go
            gopls
          ];
        };
      });

      nixosModules.default = import ./nix/nixos.nix;
      hmModules.default = import ./nix/home-manager.nix;
    };
}
