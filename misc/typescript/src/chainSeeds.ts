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
  nLeafIndex: number
): Buffer {
  proof.forEach((hash: Buffer, nIndex: number): void => {
    const leafIndex: UINT64 = UINT64(nLeafIndex.toString());
    const index: UINT64 = UINT64(nIndex.toString());
    if (leafIndex.shiftRight(index).and(UINT64(1)).equals(UINT64(0))) {
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
