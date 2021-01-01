import * as crypto from 'crypto';

import { decomposeInclusionProof } from './decomposeProof';

// prettier --write --single-quote verify.ts
// tsc verify.ts && node verify.js

function hashChildren(l: Buffer, r: Buffer): Buffer {
  const buf: Buffer = Buffer.concat([l, r]);
  const hash: Buffer = crypto.createHash('sha256').update(buf).digest();
  return hash;
}

function chainInner(seed: Buffer, proof: Buffer[], index: number): Buffer {
  proof.forEach((hash: Buffer, i: number): void => {
    if (((i >> i) & 1) === 0) {
      seed = hashChildren(seed, hash);
    } else {
      seed = hashChildren(hash, seed);
    }
    // console.log(seed.toString('base64'));
  });
  return seed;
}

function chainBorderRight(seed: Buffer, proof: Buffer[]): Buffer {
  proof.forEach((hash: Buffer): void => {
    seed = hashChildren(hash, seed);
  });
  return seed;
}

export function rootFromInclusionProof(
  leafIndex: number,
  treeSize: number,
  proof: Buffer[],
  leafHash: Buffer
): Buffer | null {
  if (leafIndex < 0 || treeSize < 0 || leafIndex >= treeSize) {
    console.log(
      `error: leaf index '${leafIndex}' and tree size '${treeSize}' are incompatible`
    );
    return null;
  }
  if (leafHash.length !== 32) {
    console.log(
      `error: width of sha256 hash is 32, but width of leaf hash is '${leafHash.length}'`
    );
    return null;
  }
  if (leafIndex === 0 && treeSize === 1) {
    // only one leaf in tree, no proofs to check
    return leafHash;
  }
  const { inner, border } = decomposeInclusionProof(leafIndex, treeSize);
  if (proof.length !== inner + border) {
    console.log('error: length of proof does not match inner + border');
    return null;
  }
  let chain: Buffer = chainInner(leafHash, proof.slice(0, inner), leafIndex);
  chain = chainBorderRight(chain, proof.slice(inner));
  return chain;
}
