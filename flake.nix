{
  description = "Keye - Key-Value DB with the ability to watch over keys";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
  outputs = { self, nixpkgs, ... }@inputs:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      version = "23.12";
    in
    {
      formatter.${system} = pkgs.nixpkgs-fmt;
      packages.${system} = rec {
        keye = pkgs.buildGoModule {
          pname = "keye";
          version = version;
          src = ./.;
          vendorHash = "sha256-8446aGvsJSE5DWqotsB5KD8hKUB877br/DI5its5hvI=";
          CGO_ENABLED = 0;
          subPackages = [ "cmd/keye" ];
        };
        dockerImage = pkgs.dockerTools.buildImage {
          name = "murtazau/keye";
          tag = version;
          config = {
            Cmd = [ "${keye}/bin/keye" ];
            WorkingDir = "/data";
          };
        };
        default = dockerImage;
      };
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          go-tools
          gopls
          protobuf
          protoc-gen-go
          protoc-gen-go-grpc
          grpcurl
        ];
      };
    };
}
