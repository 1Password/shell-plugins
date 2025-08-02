{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils = { url = "github:numtide/flake-utils"; };
  };

  outputs = inputs@{ self, nixpkgs, flake-utils, ... }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfreePredicate = pkg:
              builtins.elem (nixpkgs.lib.getName pkg) [
                "1password-cli"
              ];
          };
        };
      in
      {
        apps.supported-plugins = {
          type = "app";
          program = "${self.packages.${system}.supported-plugins}/bin/supported-plugins";
          meta = {
            description = "Generate a Nix expression containing an array of plugins that are currently supported by 1Password.";
          };
        };

        packages.supported-plugins = pkgs.writeShellApplication {
          name = "supported-plugins";
          runtimeInputs = [ pkgs._1password ];
          text = ''

            # Get the supported plugins separated by line breaks
            SUPPORTED_PLUGINS=$(op plugin list | cut -d ' ' -f1 | tail -n +2)

            if [ -z "$SUPPORTED_PLUGINS" ]; then
              echo "Error: No plugins found when calling 'op plugin list' command." >&2
              exit 1
            fi

            echo "# This file was automatically generated using 'nix run .#supported-plugins'"
            echo "$SUPPORTED_PLUGINS" | awk 'BEGIN { print "[" } {print "  \""$0"\""} END { print "]" }'

          '';
        };

        devShells.default = pkgs.mkShell {
          name = "Shell with Go toolchain";
          packages = with pkgs; [ go gopls ];
        };
      })) // {
        nixosModules.default = import ./nix/nixos.nix;
        hmModules.default = import ./nix/home-manager.nix;
      };
}
