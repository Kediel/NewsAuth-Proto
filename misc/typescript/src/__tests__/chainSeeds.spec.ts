import { hashChildren } from '../chainSeeds';

test('hashChildren', () => {
  const b64Seed: string = 'HRD7QHvXXYSEcEfshyQMWIEnavYo0gKghOGmPsur2Zk=';
  const b64Proof: string = 'uEyLjXl8SvjVS657SIjyi4bVubaE8b8Gc0n0sm5hnfw=';
  const seed: Buffer = Buffer.from(b64Seed, 'base64');
  const proof: Buffer = Buffer.from(b64Proof, 'base64');
  const b64Expected1: string = 'Sct8aVnZORh9cQYyXsNoJ9fQPnkFuaXNYLVZeYkYwpE=';
  const hashed1: Buffer = hashChildren(seed, proof);
  expect(hashed1.toString('base64')).toEqual(b64Expected1);
  const b64Expected2: string = 'Xzhq+NdbeWpF0WXWEIK/fUD/t8LX84c6wR3lspY1Geg=';
  const hashed2: Buffer = hashChildren(proof, seed);
  expect(hashed2.toString('base64')).toEqual(b64Expected2);
});
