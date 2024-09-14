{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    systems.url = "systems";
  };

  outputs =
    {
      self,
      nixpkgs,
      systems,
      ...
    }@inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {
      packages = forEachSystem (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
          };
        in
        {
          default = pkgs.buildGoModule {
            name = "http-to-aws-proxy";
            src = ./.;
            vendorHash = "sha256-mU4v2uZGOQMltpDEKJ0yKUwM1LLp5meQuJVzgMOF9Gk=";
          };
        }
      );
      devShells = forEachSystem (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
          };
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
            ];
          };
        }
      );
    };
}
