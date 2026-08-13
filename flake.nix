{
  description = "Tiny, simple, but powerful CLI framework for modern Go";

  inputs = {
    nixpkgs.url = "https://channels.nixos.org/nixpkgs-unstable/nixexprs.tar.xz";
    flake-parts.url = "github:hercules-ci/flake-parts";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      systems = [
        "aarch64-linux"
        "x86_64-linux"
        "aarch64-darwin"
      ];

      perSystem = { pkgs, ... }: {
        treefmt.programs = {
          deadnix.enable = true;
          nixfmt.enable = true;
          shellcheck.enable = true;
          shfmt.enable = true;
          statix.enable = true;
          typos.enable = true;
          yamlfmt = {
            enable = true;
            settings.formatter = {
              type = "basic";
              eof_newline = true;
              indent = 2;
              pad_line_comments = 2;
              retain_line_breaks_single = true;
              scan_folded_as_literal = true;
              trim_trailing_whitespace = true;
            };
          };
        };

        # Use golangci-lint fmt as a Go formatter
        treefmt.settings.formatter.golangci-lint-fmt = {
          command = "${pkgs.golangci-lint}/bin/golangci-lint";
          options = [ "fmt" ];
          includes = [ "*.go" ];
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            awscli2
            charm-freeze
            fd
            go
            golangci-lint
            golangci-lint-langserver
            goperf
            gopls
            gotools
            mise
            nilaway
            pkgsite
            typos
            vhs
          ];

          shellHook = ''
            echo "👋🏻 Welcome to the cli devShell!"
          '';
        };
      };
    };
}
