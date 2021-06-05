import * as crypto from 'crypto';
import * as fs from 'fs';

import { rootFromInclusionProof } from './generateRoot';

function getLogLeafHash(id: Number, postDataRaw: string): Buffer {
  function getMapLeafHash(s: string): Buffer {
    return crypto.createHash('sha256').update(s).digest();
  }
  const b64MapLeafHash: string = getMapLeafHash(postDataRaw).toString('base64');
  const logLeaf: string = String.fromCharCode(0) + id + ',' + b64MapLeafHash;
  return crypto.createHash('sha256').update(logLeaf).digest();
}

export function verify(): void {
  const rootInfoFile = process.argv[2];
  const postFile = process.argv[3];
  const proofsFile = process.argv[4];

  const rootData = JSON.parse(fs.readFileSync(rootInfoFile, 'utf8'));
  if (!rootData || !rootData.LogRoot || !rootData.LogRoot.RootHash) {
    console.log('fatal error: first argument must be name of file containing root information in json format');
    process.exit(1)
  }

  const postDataRaw = fs.readFileSync(postFile, 'utf8');
  const postData = JSON.parse(postDataRaw);
  if (!postData || !postData.ID) {
    console.log('fatal error: second argument must be name of file containing post data in json format');
    process.exit(1)
  }

  const proofData = JSON.parse(fs.readFileSync(proofsFile, 'utf8'));
  if (!proofData || !proofData.LogInclusionProof) {
    console.log('fatal error: third argument must be name of file representing proof information in json format');
    process.exit(1)
  }

  const leafHash: Buffer = getLogLeafHash(postData.ID, postDataRaw);
  const leafIndex = proofData.LogLeafIndex;
  const treeSize = proofData.LogTreeSize;
  const proofs: Buffer[] = proofData.LogInclusionProof.map((b64Proof: string): Buffer =>
    Buffer.from(b64Proof, 'base64')
  );
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proofs,
    leafHash
  );

  if (computedRoot.toString('base64') != rootData.LogRoot.RootHash) {
    console.log('false');
    process.exit(1)
  }

  console.log('true');
}

verify();
