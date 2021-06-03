# VerConf
VerConf is an internal GitHub action used to generate a program and list of target versions for the build matrix from the changed files in the current PR or merge commit.

## BuildSi Sandbox Instructions
```yaml
buildsi:
  release: 1 # Increment every time you'd like to trigger a build.
  versions:
    all:
      variants: 
        - +plugins
    2.1.4:
      variants_only: true
      variants:
      - +mpi
spack:
  specs: [abyss]
```