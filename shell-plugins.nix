{ pkgs, lib, config, is-home-manager, ... }:
with lib;
let cfg = config.programs.op-shell-plugins;
in {
  options = {
    programs.op-shell-plugins = {
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
    aliases = ''
      export OP_PLUGIN_ALIASES_SOURCED=1
      ${concatMapStrings
      (plugin: ''alias ${plugin}="op plugin run -- ${plugin}"'')
      (map (package: builtins.baseNameOf (lib.getExe package)) cfg.plugins)}
    '';
    # install the 1Password CLI, as well as any of the CLIs for which
    # shell plugins are enabled
    packages = [ pkgs._1password ] ++ cfg.plugins;
  in mkIf cfg.enable (mkMerge [
    ({
      # the option names are slightly different depending on whether you're using home-manager or not
      programs = if is-home-manager then {
        fish.interactiveShellInit = ''
          ${aliases}
        '';
        bash.initExtra = ''
          ${aliases}
        '';
        zsh.initExtra = ''
          ${aliases}
        '';
      } else {
        fish.interactiveShellInit = "";
        bash.interactiveShellInit = ''
          ${aliases}
        '';
        zsh.interactiveShellInit = ''
          ${aliases}
        '';
      };
    } // optionalAttrs is-home-manager { home.packages = packages; }
      // optionalAttrs (not is-home-manager) {
        environment.systemPackages = packages;
      })
  ]);
}
