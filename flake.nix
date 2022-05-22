{
  description = "MGA: Modern Go Application tool";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gopkgs.url = "github:sagikazarmark/go-flake";
    gopkgs.inputs.nixpkgs.follows = "nixpkgs";
    gobin.url = "github:sagikazarmark/go-bin-flake";
    gobin.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gopkgs, gobin, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;

          overlays = [
            gopkgs.overlay

            (
              final: prev: {
                golangci-lint = gobin.packages.${system}.golangci-lint-bin;
              }
            )
          ];
        };

        buildDeps = with pkgs; [ git go_1_18 gnumake ];
        devDeps = with pkgs; buildDeps ++ [
          golangci-lint
        ];
      in
      { devShell = pkgs.mkShell { buildInputs = devDeps; }; }
    );
}
