# Logo Assets

This directory contains the SVG source files for the project's logos and a `Taskfile.yml` to automate the generation of raster image formats (PNG and JPEG).

## SVG Source Files

- `logo-no-name-gradient.svg`: Logo without text, with a color gradient.
- `logo-with-name-gradient.svg`: Logo with text, with a color gradient.
- `logo-no-name-nogradient.svg`: Logo without text, with solid color segments (no gradient).
- `logo-with-name-nogradient.svg`: Logo with text, with solid color segments (no gradient).

## Image Generation Policy

The `Taskfile.yml` is used to generate various PNG and JPEG versions of these SVG logos. It's critically important to understand the characteristics and policy for these raster formats:

**WARNING: Generated PNG and JPEG files are derived from the SVG source code. They should NEVER be manually edited.** Any changes or modifications must be made to the original `.svg` files. After modifying an SVG, run `task generate` to regenerate all raster assets.

- **PNG (.png):** These files support transparency, making them suitable for use on various backgrounds. The `generate` task creates transparent PNGs directly from the SVGs.

- **JPEG (.jpeg):** These files *do not* support transparency. When converting from SVG or PNG to JPEG, a background color is applied.
    - By default, JPEGs generated directly from SVGs will have a **white background**.
    - The `generate` task also creates specific JPEG versions with a **black background** (e.g., `*-black.jpeg`) by flattening the transparent PNGs onto a black canvas.

Always choose the appropriate file format and background variant based on your specific use case and the background color where the logo will be displayed.

## Usage

To manage the generated image assets, use the following `task` commands:

- `task setup`: Checks for required tools (`svgexport` and `convert` from ImageMagick).
- `task generate`: Generates all PNG and JPEG versions of the logos from the SVG source files. This includes transparent PNGs, JPEGs with a white background, and JPEGs with a black background.
- `task clean`: Removes all generated `.png` and `.jpeg` files from this directory.
- `task`: Runs the `generate` task by default.
