import { rootFromInclusionProof } from '../generateRoot';

test('rootFromInclusionProof - only one leaf', () => {
  const b64RootHash: string = 'uEyLjXl8SvjVS657SIjyi4bVubaE8b8Gc0n0sm5hnfw=';
  const b64LeafHash: string = 'uEyLjXl8SvjVS657SIjyi4bVubaE8b8Gc0n0sm5hnfw=';
  const leafHash: Buffer = Buffer.from(b64LeafHash, 'base64');
  const leafIndex = 0;
  const treeSize = 1;
  const proofs: Buffer[] = [];
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proofs,
    leafHash
  );
  expect(computedRoot.toString('base64')).toEqual(b64RootHash);
});

test.only('rootFromInclusionProof - two leaves', () => {
  const b64RootHash: string = '0hcIpiYpMMUsS/lKLMlYzswE7vU6JDOq6Wm5QmbEcYg=';
  const b64LeafHash: string = 'HRD7QHvXXYSEcEfshyQMWIEnavYo0gKghOGmPsur2Zk=';
  const b64Proofs: string[] = ['uEyLjXl8SvjVS657SIjyi4bVubaE8b8Gc0n0sm5hnfw='];
  const leafHash: Buffer = Buffer.from(b64LeafHash, 'base64');
  const leafIndex = 1;
  const treeSize = 2;
  const proofs: Buffer[] = b64Proofs.map((b64Proof: string): Buffer =>
    Buffer.from(b64Proof, 'base64')
  );
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proofs,
    leafHash
  );
  expect(computedRoot.toString('base64')).toEqual(b64RootHash);
});
