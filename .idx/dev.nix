{ pkgs, ... }:

{ # The dev.nix file is the entrypoint of your environment.
  # It contains all the packages and configurations that you need
  # for your project.
  channel = "stable-24.05"; # or "unstable"

  # A list of packages that are needed for your development environment.
  packages = [
    pkgs.go-outline
    pkgs.gopkgs
    pkgs.gopls
    pkgs.gotools
    pkgs.go
    pkgs.tesseract
    pkgs.pkg-config
    pkgs.haskellPackages.snap-templates
  ];

  # Sets environment variables in the workspace
  env = {};
  idx = {
    # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
    extensions = [
      # "vscodevim.vim"
      "golang.go"
    ];

    # Enable previews
    previews = {
      enable = true;
      previews = {
        # web = {
        #   # Example: run "npm run dev" with PORT set to IDX's defined port for previews,
        #   # and show it in IDX's web preview panel
        #   command = ["npm" "run" "dev"];
        #   manager = "web";
        #   env = {
        #     # Environment variables to set for your server
        #     PORT = "$PORT";
        #   };
        # };
      };
    };

    # Workspace lifecycle hooks
    workspace = {
      # Runs when a workspace is first created
      onCreate = {
        # Example: install JS dependencies from NPM
        # npm-install = "npm install";
      };
      # Runs when the workspace is (re)started
      onStart = {
        # Example: start a background task to watch and re-build backend code
        # watch-backend = "npm run watch-backend";
      };
    };
  };
}
