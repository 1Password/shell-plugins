{
  pkgs,
  lib,
  config,
  is-home-manager,
  ...
}:
with lib;
let
  cfg = config.programs._1password-shell-plugins;

  opPkg = if cfg.package != null then cfg.package else pkgs._1password-cli;

  getExeName =
    package:
    # NOTE: SAFETY: This is okay because the `packages` list is also referred
    # to below as `home.packages = packages;` or `environment.systemPackages = packages;`
    # depending on if it's using `home-manager` or not; this means that Nix can still
    # compute the dependency tree, even though we're discarding string context here,
    # since the packages are still referred to below without discarding string context.
    strings.unsafeDiscardStringContext (baseNameOf (getExe package));

  mkPluginSupportCheck =
    pluginExeNames:
    let
      configuredFile = pkgs.writeText "op-configured-plugins.txt" (
        # ensure trailing newline
        lib.concatStringsSep "\n" pluginExeNames + "\n"
      );
    in
    pkgs.runCommand "op-shell-plugins-support-check"
      {
        nativeBuildInputs = [
          pkgs.coreutils
          pkgs.gnugrep
          pkgs.gawk
          pkgs.gnused
        ];
      }
      ''
        		set -euo pipefail

        		export XDG_CONFIG_HOME="$TMPDIR/xdg-config"
        		mkdir -p "$XDG_CONFIG_HOME"

        		"${opPkg}/bin/op" plugin list \
        		  | awk 'NR>1 { print $1 }' \
        		  | sed '/^$/d' \
        		  | sort -u > supported.txt

        		if [ ! -s supported.txt ]; then
        		  echo "ERROR: \`op plugin list\` produced no supported plugins (unexpected output or CLI failure)." >&2
        		  echo "Raw output was:" >&2
        		  "${opPkg}/bin/op" plugin list >&2 || true
        		  exit 1
        		fi

        		missing=0
        		while IFS= read -r plugin; do
        		  [ -z "$plugin" ] && continue
        		  if ! grep -Fxq "$plugin" supported.txt; then
        		    echo "ERROR: Configured plugin '$plugin' is not supported by this op binary (\`${opPkg.name or "op"}\`)." >&2
        		    missing=1
        		  fi
        		done < "${configuredFile}"

        		if [ "$missing" -ne 0 ]; then
        		  echo "" >&2
        		  echo "Supported plugins according to \`op plugin list\`:" >&2
        		  cat supported.txt >&2
        		  exit 1
        		fi
        		mkdir -p "$out"
        		echo "Plugin support check passed." > "$out/check.txt"

        				  '';

in
{
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
        description = "CLI Packages to enable 1Password Shell Plugins for; ensure that a Shell Plugin exists by checking the docs: https://developer.1password.com/docs/cli/shell-plugins/";

      };
    };
  };

  config =
    let

      # executable names as strings, e.g. `pkgs.gh` => `"gh"`, `pkgs.awscli2` => `"aws"`
      pkg-exe-names = map getExeName cfg.plugins;
      plugin-support-check = mkPluginSupportCheck pkg-exe-names;
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
      initExtraPosix = strings.concatStringsSep "\n" posixFunctions;
    in
    mkIf cfg.enable (mkMerge [
      {
        programs.fish.interactiveShellInit = strings.concatStringsSep "\n" fishFunctions;
      }
      (optionalAttrs is-home-manager {
        home.checks = [ plugin-support-check ];
        programs = {
          # for the Bash and Zsh home-manager modules,
          # the initExtra/initContent option is equivalent to Fish's interactiveShellInit
          bash.initExtra = initExtraPosix;
          zsh.initContent = initExtraPosix;
        };
        home = {
          inherit packages;
          sessionVariables = {
            OP_PLUGINS_SOURCED = "1";
          };
        };
      })
      (optionalAttrs (!is-home-manager) {
        system.checks = [ plugin-support-check ];
        programs = {
          bash.interactiveShellInit = strings.concatStringsSep "\n" posixFunctions;
          zsh.interactiveShellInit = strings.concatStringsSep "\n" posixFunctions;
        };
        environment = {
          systemPackages = packages;
          variables = {
            OP_PLUGINS_SOURCED = "1";
          };
        };
      })
    ]);
}
