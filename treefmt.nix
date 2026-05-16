{ pkgs, ... }:
{
  projectRootFile = "flake.nix";

  programs = {
    # keep-sorted start block=yes
    deadnix.enable = true;
    keep-sorted.enable = true;
    nixfmt = {
      enable = true;
      package = pkgs.nixfmt-rfc-style;
      includes = [ "*.nix" ];
      strict = true;
    };
    statix.enable = true;
    yamlfmt = {
      enable = true;
      settings = {
        formatter = {
          eof_newline = true;
          max_line_length = -1;
          retain_line_breaks_single = true;
          scan_folded_as_literal = true;
          trim_trailing_whitespace = true;
        };
      };
    };
    # keep-sorted end
  };

  settings = {
    global.excludes = [
      "*.lock"
      "*.patch"
      "LICENSE*"
      ".github/**/*.md"
    ];
  };

  settings.formatter = {
    keep-sorted = {
      options = [
        "--mode"
        "fix"
      ];
    };
  };
}
