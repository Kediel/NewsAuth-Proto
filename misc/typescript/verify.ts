import * as crypto from "crypto";

// prettier --write --single-quote verify.ts
// tsc verify.ts && node verify.js

function decomposeInclusionProof(
  index: number,
  size: number
): { innerSize: number; borderSize: number } {
  function getInnerProofSize(index: number, size: number): number {
    function getMinNumBitsNeeded(num: number): number {
      // Get number of bits needed to represent positive integer:
      // https://stackoverflow.com/questions/12349498/minimum-number-of-bits-to-represent-number
      return Math.Floor(Math.Log(num, 2)) + 1;
    }
    return getMinNumBitsNeeded(index ^ (size - 1));
  }

  function getOnesCount(n: number): number {
    // https://stackoverflow.com/questions/43122082/efficiently-count-the-number-of-bits-in-an-integer-in-javascript
    n = n - ((n >> 1) & 0x55555555);
    n = (n & 0x33333333) + ((n >> 2) & 0x33333333);
    return (((n + (n >> 4)) & 0xf0f0f0f) * 0x1010101) >> 24;
  }

  innerSize = getInnerProofSize(index, size);
  borderSize = getOnesCount(index >> ToUint32(inner));
  return { innerSize, borderSize };
}

function rootFromInclusionProof(
  leafIndex: number,
  treeSize: number,
  proof: Buffer[],
  leafHash: Buffer
): Buffer | null {
  if (leafIndex < 0 || treeSize < 0 || leafIndex >= treeSize) return null;
  // TODO: check that length is compatible with hash algorithm

  const { innerSize, borderSize } = decomposeInclusionProof(index, size);
  if (proof.length !== innerSize + borderSize) {
    return null;
  }
  return Buffer.from("cats");
}

// const base64Proofs: string[] = [
//   "swDUQGPbJoO2TQg7rdZad0Smv2QNcFa3FWvaFGCfaoI=",
//   "Yfjth1PvCCwnz6DiesNKiNyxvKuitzFWxRoJFllFCgg=",
//   "hQRlPYRSvSoRNwDxDs+B8Gfu1RfctUjP2H060cLJpgI=",
//   "i225NGeEZiK8rhgkEeq/ryxKXALHmpj280oSK61LfG8=",
//   "XEYlWuPyvUvyzS4s6BPXIEcLQm2WaX5m5BpEganXSfs="
// ];
//
// const proofs: Uint8Array[] = base64Proofs.map(base64Hash =>
//   Uint8Array.from(Buffer.from(base64Hash, "base64").toString(), c =>
//     c.charCodeAt(0)
//   )
// );
// const b: Buffer = Buffer.concat([proofs[0], proofs[1]]);
//
// const h = crypto.createHash("sha256").update(b).digest("base64");
//
// console.log(h, base64Proofs[2]);
