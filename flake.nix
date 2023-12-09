{
  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2311.552856.tar.gz";
    gitignore.url = "github:hercules-ci/gitignore.nix";
  };

  outputs = { self, nixpkgs, gitignore }: 
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in rec {
      # development shell
      devShells.${system}.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
          ];
      };

      # build the go application
      application = pkgs.buildGoPackage {
        name = "gazes-proxy";
        src = gitignore.lib.gitignoreSource ./.;
        goPackagePath = "github.com/trail-l31/gazes-proxy";
      };

      # build the docker image
      image = pkgs.dockerTools.buildImage rec {
        name = "gazes-proxy";
        tag = "latest";
        created = "now";

        contents = application;

        config = {
          Cmd = [ "${application}/bin/cmd" ];
          ExposedPorts = {
            "3000/tcp" = {};
          };
        };
      };
    };
}
