# Conda Forge Recipe for ensor

This directory contains the conda-forge recipe for building the ensor package using the modern CEP-13 (v1 / rattler-build) format.

## Building Locally

```bash
# Install dependencies
pixi install

# Build the conda package with rattler-build
pixi run conda-build
```

## Publishing to Conda Forge

To publish `ensor` to conda-forge for the first time:

1. **Tag and Push the release** on GitHub:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. **Compute the SHA256 checksum** of the source archive:
   ```bash
   curl -sL https://github.com/Pratyay360/ensor/archive/refs/tags/v0.1.0.tar.gz | sha256sum
   ```
   Update the `sha256` field in `recipe/recipe.yaml`.

3. **Fork and Clone `conda-forge/staged-recipes`**:
   ```bash
   git clone --depth=1 https://github.com/<your-username>/staged-recipes.git
   cd staged-recipes
   git checkout -b add-ensor
   ```

4. **Copy the recipe**:
   ```bash
   mkdir -p recipes/ensor
   cp /path/to/ensor/recipe/recipe.yaml recipes/ensor/
   ```

5. **Commit and open a Pull Request**:
   ```bash
   git add recipes/ensor/recipe.yaml
   git commit -m "Add recipe for ensor"
   git push origin add-ensor
   ```
   Open a PR against `conda-forge/staged-recipes`.

Once the PR passes CI checks and is merged by the conda-forge team:
- A new repository `conda-forge/ensor-feedstock` is automatically created.
- Packages are built and published to the `conda-forge` channel.
- Future versions are updated by submitting PRs to `conda-forge/ensor-feedstock` (or via conda-forge's automated bot).
