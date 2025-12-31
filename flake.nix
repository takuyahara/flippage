{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        formatter = pkgs.nixfmt-tree;
        devShells = {
          default = pkgs.mkShell {
            name = "flippage";
            packages = with pkgs; [
              jq
              gitflow
              bash-completion
              pinact
              zizmor
              # Go tools
              go
              gotools
              go-tools
              gopls
              go-outline
              gopkgs
              gocode-gomod
              godef
            ];
            shellHook = ''
              if [[ -e ./.vscode/settings.json ]]; then
                goroot="${pkgs.go}/share/go"
                gopls="${pkgs.gopls}/bin/gopls"
                dlv="${pkgs.delve}/bin/dlv"
                staticcheck="${pkgs.go-tools}/bin/staticcheck"
                cat <<< $(cat .vscode/settings.json | \
                  jq ".[\"go.goroot\"] = \"$goroot\"" | \
                  jq ".[\"go.alternateTools\"].gopls = \"$gopls\"" | \
                  jq ".[\"go.alternateTools\"].dlv = \"$dlv\"" | \
                  jq ".[\"go.alternateTools\"].staticcheck = \"$staticcheck\"" \
                ) > .vscode/settings.json
              fi
              . "${pkgs.bash-completion}/etc/profile.d/bash_completion.sh"
              PATH=$PATH:~/go/bin
            '';
          };
        };
      }
    );
}
