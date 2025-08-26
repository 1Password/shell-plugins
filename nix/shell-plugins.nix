{ pkgs, lib, config, is-home-manager, ... }:
with lib;
let
  cfg = config.programs._1password-shell-plugins;

  supported_plugins = splitString "\n" (lib.readFile "${
    # get the list of supported plugin executable names
      pkgs.runCommand "op-plugin-list" { }
      # 1Password CLI tries to create the config directory automatically, so set a temp XDG_CONFIG_HOME
      # since we don't actually need it for this
      "mkdir $out && XDG_CONFIG_HOME=$out ${pkgs._1password}/bin/op plugin list | cut -d ' ' -f1 | tail -n +2 > $out/plugins.txt"
    }/plugins.txt");
  getExeName = package:
    # NOTE: SAFETY: This is okay because the `packages` list is also referred
    # to below as `home.packages = packages;` or `environment.systemPackages = packages;`
    # depending on if it's using `home-manager` or not; this means that Nix can still
    # compute the dependency tree, even though we're discarding string context here,
    # since the packages are still referred to below without discarding string context.
    strings.unsafeDiscardStringContext (baseNameOf (getExe package));
in {
  options = {
    programs._1password-shell-plugins = {
      enable = mkEnableOption "1Password Shell Plugins";
      package = mkPackageOption pkgs "_1password-cli" { nullable = true; };
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
        # this is a bit of a hack to do option validation;
        # ensure that the list of packages include only packages
        # for which the executable has a supported 1Password Shell Plugin
        apply = package_list:
          map (package:
            if (elem (getExeName package) supported_plugins) then
              package
            else
              abort "${
                getExeName package
              } is not a valid 1Password Shell Plugin. A list of supported plugins can be found by running `op plugin list` or at: https://developer.1password.com/docs/cli/shell-plugins/")
          package_list;
      };
    };
  };

  config = let
    # executable names as strings, e.g. `pkgs.gh` => `"gh"`, `pkgs.awscli2` => `"aws"`
    pkg-exe-names = map getExeName cfg.plugins;
    # Explanation:
    # Map over `cfg.plugins` (the value of the `plugins` option provided by the user)
    # and for each package specified, get the executable name, then create a shell function
    # of the form:
    #
    # For Bash and Zsh:
    # ```
    #   {pkg}() {
    #     op plugin run -- {pkg};
    #   }
    # ```
    #
    # And for Fish:
    # ```
    #  function {pkg} --wraps {pkg}
    #    op plugin run -- {pkg}
    #  end
    # ```
    # where `{pkg}` is the executable name of the package
    posixFunctions = map (package: ''
      ${package}() {
        op plugin run -- ${package} "$@";
      }
    '') pkg-exe-names;
    fishFunctions = map (package: ''
      function ${package} --wraps "${package}" --description "1Password Shell Plugin for ${package}"
        op plugin run -- ${package} $argv
      end
    '') pkg-exe-names;
    packages = lib.optional (cfg.package != null) cfg.package ++ cfg.plugins;
  in mkIf cfg.enable (mkMerge [
    (
      optionalAttrs is-home-manager {
      programs = {
        # for the Bash and Zsh home-manager modules,
        # the initExtra option is equivalent to Fish's interactiveShellInit
        fish.interactiveShellInit = strings.concatStringsSep "\n" fishFunctions;
        bash.initExtra = strings.concatStringsSep "\n" posixFunctions;
        zsh.initExtra = strings.concatStringsSep "\n" posixFunctions;
      };
      home = {
        inherit packages;
        sessionVariables = { OP_PLUGINS_SOURCED = "1"; };
      };
    } // optionalAttrs (!is-home-manager) {
      programs = {
        bash.interactiveShellInit =
          strings.concatStringsSep "\n" posixFunctions;
        zsh.interactiveShellInit = strings.concatStringsSep "\n" posixFunctions;
      };
      environment = {
        systemPackages = packages;
        variables = { OP_PLUGINS_SOURCED = "1"; };
      };
    })
  ]);
}
