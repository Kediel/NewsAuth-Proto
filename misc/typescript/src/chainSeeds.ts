import { UINT64 } from 'cuint';
import * as crypto from 'crypto';

export function hashChildren(l: Buffer, r: Buffer): Buffer {
  const RFC6962NodeHashPrefix: Buffer = Buffer.alloc(1, 1);
  const buf: Buffer = Buffer.concat([RFC6962NodeHashPrefix, l, r]);
  const hash: Buffer = crypto.createHash('sha256').update(buf).digest();
  return hash;
}

export function chainInner(
  seed: Buffer,
  proof: Buffer[],
  leafIndex: number
): Buffer {
  const li: UINT64 = UINT64(leafIndex.toString());
  proof.forEach((hash: Buffer, index: number): void => {
    const i: UINT64 = UINT64(index.toString());
    if (li.shiftRight(i).clone().and(UINT64(1)).equals(UINT64(0))) {
      seed = hashChildren(seed, hash);
    } else {
      seed = hashChildren(hash, seed);
    }
  });
  return seed;
}

export function chainBorderRight(seed: Buffer, proof: Buffer[]): Buffer {
  proof.forEach((hash: Buffer): void => {
    seed = hashChildren(hash, seed);
  });
  return seed;
}
