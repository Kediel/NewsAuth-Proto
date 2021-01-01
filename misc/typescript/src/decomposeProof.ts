import { UINT64 } from 'cuint';

export function numToUINT64(n: number): UINT64 {
  return UINT64(n.toString());
}

export function getMinNumBits(n: UINT64): number {
  const zero: UINT64 = UINT64(0);
  if (n.or(zero).equals(zero)) return 0; // n is zero
  let mask = UINT64(1);
  let bitsNeeded = 1;
  while (mask.lessThan(n) === true) {
    mask = mask.shiftLeft(1).or(UINT64(1));
    bitsNeeded++;
  }
  return bitsNeeded;
}

export function getNumSetBits(n: UINT64): number {
  const zero: UINT64 = UINT64(0);
  if (n.or(zero).equals(zero)) return 0; // n is zero
  let bits = 0;
  let mask = UINT64(1);
  for (let i = 0; i < 64; i++) {
    if (n.clone().and(mask).equals(mask)) bits++;
    mask.shiftLeft(1);
  }
  return bits;
}

export function decomposeInclusionProof(
  index: number,
  size: number
): { inner: number; border: number } {
  const i: UINT64 = numToUINT64(index);
  const s: UINT64 = numToUINT64(size);
  const t: UINT64 = i.clone().xor(s.subtract(numToUINT64(1)));
  const inner: number = getMinNumBits(t);
  const border: number = getNumSetBits(i.shiftRight(numToUINT64(inner)));
  return { inner, border };
}
