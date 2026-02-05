#!/usr/bin/env node
const fs = require('fs').promises;
const path = require('path');
const sharp = require('sharp');

const SRC_DIRS = [
  path.join(__dirname, '..', 'static', 'img'),
  path.join(__dirname, '..', 'static', 'images')
];

const SIZES = [320, 640, 1024, 1600];

function isImage(filename) {
  return /\.(jpe?g|png)$/i.test(filename);
}

async function walk(dir) {
  let results = [];
  try {
    const entries = await fs.readdir(dir, { withFileTypes: true });
    for (const entry of entries) {
      const full = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        results = results.concat(await walk(full));
      } else if (entry.isFile() && isImage(entry.name)) {
        results.push(full);
      }
    }
  } catch (err) {
    return results;
  }
  return results;
}

async function ensureDir(filePath) {
  const dir = path.dirname(filePath);
  await fs.mkdir(dir, { recursive: true });
}

async function processFile(file) {
  try {
    const image = sharp(file);
    // Generate size variants for webp and avif
    for (const w of SIZES) {
      const ext = path.extname(file);
      const base = file.slice(0, -ext.length);
      const webp = `${base}-${w}.webp`;
      const avif = `${base}-${w}.avif`;
      await ensureDir(webp);
      await image.clone().resize({ width: w }).webp({ quality: 80 }).toFile(webp);
      await image.clone().resize({ width: w }).avif({ quality: 50 }).toFile(avif);
    }
    // Also write full-size fallback webp/avif
    const fullWebp = file.replace(/\.(jpe?g|png)$/i, '.webp');
    const fullAvif = file.replace(/\.(jpe?g|png)$/i, '.avif');
    await image.clone().webp({ quality: 80 }).toFile(fullWebp);
    await image.clone().avif({ quality: 50 }).toFile(fullAvif);
    console.log('Converted:', path.relative(process.cwd(), file));
  } catch (err) {
    console.error('Error processing', file, err.message);
  }
}

async function main() {
  const files = new Set();
  for (const dir of SRC_DIRS) {
    const found = await walk(dir);
    found.forEach(f => files.add(f));
  }
  if (files.size === 0) {
    console.log('No source images found in static/img or static/images.');
    return;
  }
  for (const f of files) await processFile(f);
}

main().catch(err => {
  console.error(err);
  process.exit(1);
});
