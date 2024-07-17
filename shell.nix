let pkgs = import <nixpkgs> {};

in pkgs.mkShell rec {
  name = "flippage";

  buildInputs = with pkgs; [
    gitflow
    bash-completion
    go
    gotools
    go-tools
    gopls
    go-outline
    gopkgs
    gocode-gomod
    godef
    darwin.apple_sdk.frameworks.Cocoa
  ];
  shellHook =
  ''
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
}
