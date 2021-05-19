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

test('rootFromInclusionProof - two leaves', () => {
  const b64RootHash: string = '0hcIpiYpMMUsS/lKLMlYzswE7vU6JDOq6Wm5QmbEcYg=';
  const b64LeafHash: string = '5OEffD4ILe07w1BKbLj0SYf6XecfJgFx3xwI+Zgyxos=';
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

test('rootFromInclusionProof - three leaves', () => {
  const b64RootHash: string = '2CKAUML+Ego8w3ZjVPg4y8v5K0wGidsc3xZ6KZ6UfH4=';
  const b64LeafHash: string = 'yTNCFGyaQnNr3iAq75Ogorz8LoAZAYqYlHKGUCi45V0=';
  const b64Proof: string[] = ['0hcIpiYpMMUsS/lKLMlYzswE7vU6JDOq6Wm5QmbEcYg='];
  const leafHash: Buffer = Buffer.from(b64LeafHash, 'base64');
  const leafIndex = 2;
  const treeSize = 3;
  const proof: Buffer[] = b64Proof.map((b64Hash: string): Buffer =>
    Buffer.from(b64Hash, 'base64')
  );
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proof,
    leafHash
  );
  expect(computedRoot.toString('base64')).toEqual(b64RootHash);
});

test('rootFromInclusionProof - 16 leaves', () => {
  const b64RootHash: string = '48WwwFM+VmnsFOhqKKdZtFLwlNtO9QS7ykivP2Gfjdw=';
  const b64LeafHash: string = 'Iat5/ZZGiFhqeYbFzMw+5QEnCtWu2rhSJwzTHFwJPM4=';
  const b64Proof: string[] = [
    'KPLw8j4oDO4f6hHUw0rhAbqeOb1zu2g5/SV+hOCQeGw=',
    'f6h5C/IhHu/cnE6xt+/4TPcZJf7/Y/jfyr20Y7DEC90=',
    'GrjCOOtxpqOkLAhJjD2XntqqzJUmrNvgfLXpU0OxYQk=',
    '6JOfmO/E62fFZTt5u7HomIdUxyjHlLOa5yIH0EmfuOY='
  ];
  const leafHash: Buffer = Buffer.from(b64LeafHash, 'base64');
  const leafIndex = 15;
  const treeSize = 16;
  const proof: Buffer[] = b64Proof.map((b64Hash: string): Buffer =>
    Buffer.from(b64Hash, 'base64')
  );
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proof,
    leafHash
  );
  expect(computedRoot.toString('base64')).toEqual(b64RootHash);
});

test('rootFromInclusionProof - 16 older leaf', () => {
  const b64RootHash: string = '48WwwFM+VmnsFOhqKKdZtFLwlNtO9QS7ykivP2Gfjdw=';
  const b64LeafHash: string = 'j0soaRyiUehPSR+PAEUrVeL+SCMuKZP3x9Q71eMuXMI=';
  const b64Proof: string[] = [
    'yTNCFGyaQnNr3iAq75Ogorz8LoAZAYqYlHKGUCi45V0=',
    '0hcIpiYpMMUsS/lKLMlYzswE7vU6JDOq6Wm5QmbEcYg=',
    'Y7mT6EZ73vAWOhG+J/Y/PmVVs32hXAUGZw0YL+onhFw=',
    'VvL6DPvcMP+EqFxhiLAG6nxYzZXsAsTFwULSpQM/4x4='
  ];
  const leafHash: Buffer = Buffer.from(b64LeafHash, 'base64');
  const leafIndex = 3;
  const treeSize = 16;
  const proof: Buffer[] = b64Proof.map((b64Hash: string): Buffer =>
    Buffer.from(b64Hash, 'base64')
  );
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proof,
    leafHash
  );
  expect(computedRoot.toString('base64')).toEqual(b64RootHash);
});

test('rootFromInclusionProof - 17', () => {
  const b64RootHash: string = 'e9lYo8MWUmjpjP75pYTq8Oo2k4NMtScn4y4MUb4LCH0=';
  const b64LeafHash: string = 'h8p9O7FMp1WhDCvrlhANACSn5ukgMESAyb63nPqCqyo=';
  const b64Proof: string[] = ['48WwwFM+VmnsFOhqKKdZtFLwlNtO9QS7ykivP2Gfjdw='];
  const leafHash: Buffer = Buffer.from(b64LeafHash, 'base64');
  const leafIndex = 16;
  const treeSize = 17;
  const proof: Buffer[] = b64Proof.map((b64Hash: string): Buffer =>
    Buffer.from(b64Hash, 'base64')
  );
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proof,
    leafHash
  );
  expect(computedRoot.toString('base64')).toEqual(b64RootHash);
});

test('rootFromInclusionProof - 17 older leaf', () => {
  const b64RootHash: string = 'e9lYo8MWUmjpjP75pYTq8Oo2k4NMtScn4y4MUb4LCH0=';
  const b64LeafHash: string = 'KPLw8j4oDO4f6hHUw0rhAbqeOb1zu2g5/SV+hOCQeGw=';
  const b64Proof: string[] = [
    'Iat5/ZZGiFhqeYbFzMw+5QEnCtWu2rhSJwzTHFwJPM4=',
    'f6h5C/IhHu/cnE6xt+/4TPcZJf7/Y/jfyr20Y7DEC90=',
    'GrjCOOtxpqOkLAhJjD2XntqqzJUmrNvgfLXpU0OxYQk=',
    '6JOfmO/E62fFZTt5u7HomIdUxyjHlLOa5yIH0EmfuOY=',
    'h8p9O7FMp1WhDCvrlhANACSn5ukgMESAyb63nPqCqyo='
  ];
  const leafHash: Buffer = Buffer.from(b64LeafHash, 'base64');
  const leafIndex = 14;
  const treeSize = 17;
  const proof: Buffer[] = b64Proof.map((b64Hash: string): Buffer =>
    Buffer.from(b64Hash, 'base64')
  );
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proof,
    leafHash
  );
  expect(computedRoot.toString('base64')).toEqual(b64RootHash);
});

test('live', () => {
  const b64RootHash: string = 'g3FnfguCfgttsfEssCRtA8drAnD+rTeMYcj\/4n9ulKc=';
  const b64LeafHash: string = 'dbBmhEytFh+sBNAFCMrl2ycF157r36BN3goLdxhgDss=';
  const b64Proof: string[] = [
    "OaOVTXjirhUfVXvlzmgTlzy2tPzgiVNN+UpiCI3NoGM=",
    "rp2KkeMnF8nF23sR//bcXrV4QoqsZHYvRccldKXzJZM=",
    "uewJNjobgH9PoEj/JyErmGsP+MY9ZdkCsPPdWVj8xuw=",
    "D2KCChpM546Wv9efhaxHoswpGaD2TwD17J2YScimfgk=",
    "uu65oAfui8iQ6SXb5TTq+Yw3KNA1OSjEXmMTAwBenio=",
    "yhjiiNzqUsgm/POz/uoFiR3MqTAyerZUrRekBFtHPVc=",
    "gIISiTanq+o48XY5RNC1AuL5YMrt9PSVEVL9vk5nzjw=",
    "ckyyc8VhBBLzgbjtnNznkIjx6FY7Ua/Qycigbw8m3rw="
  ];
  const leafHash: Buffer = Buffer.from(b64LeafHash, 'base64');
  const leafIndex = 224;
  const treeSize = 225;
  const proof: Buffer[] = b64Proof.map((b64Hash: string): Buffer =>
    Buffer.from(b64Hash, 'base64')
  );
  const computedRoot: Buffer = rootFromInclusionProof(
    leafIndex,
    treeSize,
    proof,
    leafHash
  );
  expect(computedRoot.toString('base64')).toEqual(b64RootHash);
});
