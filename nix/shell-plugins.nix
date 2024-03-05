{ pkgs, lib, config, is-home-manager, ... }:
with lib;
let cfg = config.programs._1password-shell-plugins;
in {
  options = {
    programs._1password-shell-plugins = {
      enable = mkEnableOption "1Password Shell Plugins";
      plugins = mkOption {
        type = types.listOf types.package;
        default = [ ];
        example = literalExpression ''
          with pkgs; [
            gh
            awscli2
            cachix
          ]
        '';
        description =
          "CLI Packages to enable 1Password Shell Plugins for; ensure that a Shell Plugin exists by checking the docs: https://developer.1password.com/docs/cli/shell-plugins/";
      };
    };
  };

  config = let
    # Explanation:
    # Map over `cfg.plugins` (the value of the `plugins` option provided by the user)
    # and for each package specified, get the executable name, then create a shell alias
    # of the form:
    # `alias {pkg}="op plugin run -- {pkg}"`
    # where `{pkg}` is the executable name of the package
    aliases = listToAttrs (map (package: {
      name = package;
      value = "op plugin run -- ${package}";
    }) (map (package:
      # NOTE: SAFETY: This is okay because the `packages` list is also referred
      # to below as `home.packages = packages;` or `environment.systemPackages = packages;`
      # depending on if it's using `home-manager` or not; this means that Nix can still
      # compute the dependency tree, even though we're discarding string context here,
      # since the packages are still referred to below without discarding string context.
      builtins.unsafeDiscardStringContext (baseNameOf (getExe package)))
      cfg.plugins));
    packages = [ pkgs._1password ] ++ cfg.plugins;
  in mkIf cfg.enable (mkMerge [
    ({
      programs = {
        bash.shellAliases = aliases;
        zsh.shellAliases = aliases;
        fish.shellAliases = aliases;
      };
    } // optionalAttrs is-home-manager { home.packages = packages; }
      // optionalAttrs (!is-home-manager) {
        environment.systemPackages = packages;
      })
  ]);
}
