#!/usr/bin/env node

const fs = require("node:fs");
const fsp = require("node:fs/promises");
const os = require("node:os");
const path = require("node:path");
const https = require("node:https");
const tar = require("tar");
const extractZip = require("extract-zip");

const packageJson = require("../package.json");

const OWNER = "CoderSerio";
const REPO = "pokemand-go";
const VERSION = packageJson.version;
const BIN_DIR = path.join(__dirname, "..", ".pkmg-bin");
const BIN_NAME = process.platform === "win32" ? "pkmg.exe" : "pkmg";
const DEST_BIN = path.join(BIN_DIR, BIN_NAME);

function resolveAssetName() {
  const platformMap = {
    darwin: "Darwin",
    linux: "Linux",
    win32: "Windows",
  };

  const archMap = {
    x64: "x86_64",
  };

  const osName = platformMap[process.platform];
  const archName = archMap[process.arch];

  if (!osName || !archName) {
    throw new Error(
      `Unsupported platform for current npm wrapper: ${process.platform}/${process.arch}. ` +
      "This wrapper currently supports x64 release assets only."
    );
  }

  const extension = process.platform === "win32" ? "zip" : "tar.gz";
  return `pokemand-go-${osName}-${archName}.${extension}`;
}

function download(url, destination) {
  return new Promise((resolve, reject) => {
    const request = https.get(url, (response) => {
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        response.resume();
        download(response.headers.location, destination).then(resolve).catch(reject);
        return;
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Download failed: ${response.statusCode} ${response.statusMessage}`));
        return;
      }

      const file = fs.createWriteStream(destination);
      response.pipe(file);
      file.on("finish", () => file.close(resolve));
      file.on("error", reject);
    });

    request.on("error", reject);
  });
}

async function findBinary(dir, name) {
  const entries = await fsp.readdir(dir, { withFileTypes: true });
  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      const nested = await findBinary(fullPath, name);
      if (nested) return nested;
      continue;
    }
    if (entry.name === name) {
      return fullPath;
    }
  }
  return null;
}

async function main() {
  if (process.env.PKMG_NPM_SKIP_DOWNLOAD === "1") {
    console.log("Skipping pkmg binary download because PKMG_NPM_SKIP_DOWNLOAD=1");
    return;
  }

  const assetName = resolveAssetName();
  const downloadURL = `https://github.com/${OWNER}/${REPO}/releases/download/v${VERSION}/${assetName}`;
  const tempDir = await fsp.mkdtemp(path.join(os.tmpdir(), "pkmg-npm-"));
  const archivePath = path.join(tempDir, assetName);
  const extractDir = path.join(tempDir, "extract");

  await fsp.mkdir(extractDir, { recursive: true });
  await fsp.mkdir(BIN_DIR, { recursive: true });

  console.log(`Downloading pkmg ${VERSION} binary from GitHub Release...`);
  await download(downloadURL, archivePath);

  if (assetName.endsWith(".zip")) {
    await extractZip(archivePath, { dir: extractDir });
  } else {
    await tar.x({
      file: archivePath,
      cwd: extractDir,
      gzip: true,
    });
  }

  const extractedBinary = await findBinary(extractDir, BIN_NAME);
  if (!extractedBinary) {
    throw new Error(`Could not find ${BIN_NAME} in downloaded archive ${assetName}`);
  }

  await fsp.copyFile(extractedBinary, DEST_BIN);
  if (process.platform !== "win32") {
    await fsp.chmod(DEST_BIN, 0o755);
  }

  console.log(`Installed pkmg binary to ${DEST_BIN}`);
}

main().catch((error) => {
  console.error(error.message || error);
  process.exit(1);
});
