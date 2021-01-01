import { UINT64 } from 'cuint';

import {
  decomposeInclusionProof,
  getNumSetBits,
  getMinNumBits,
  numToUINT64
} from '../decomposeProof';

test.skip('numToUINT64', () => {
  const max = Number.MAX_SAFE_INTEGER; // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
  let i = 0;
  while (i <= max) {
    expect(i.toString()).toEqual(numToUINT64(i).toString(10));
    console.log(i);
    i++;
  }
});

test('getMinNumBits', () => {
  // Testing output against: https://golang.org/pkg/math/bits/#Len64
  const basicRubric: number[] = [
    0,
    1,
    2,
    2,
    3,
    3,
    3,
    3,
    4,
    4,
    4,
    4,
    4,
    4,
    4,
    4,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    5,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    6,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7,
    7
  ];
  const computed: number[] = [];
  for (let i = 0; i < 100; i++) {
    computed.push(getMinNumBits(numToUINT64(i)));
  }
  expect(computed).toEqual(basicRubric);
  expect(getMinNumBits(numToUINT64(1024))).toEqual(11);
  expect(getMinNumBits(numToUINT64(2047))).toEqual(11);
  expect(getMinNumBits(numToUINT64(65535))).toEqual(16);
  expect(getMinNumBits(numToUINT64(65536))).toEqual(17);
  expect(getMinNumBits(numToUINT64(4294967295))).toEqual(32);
  expect(getMinNumBits(numToUINT64(4294967296))).toEqual(33);
  expect(getMinNumBits(numToUINT64(9007199254740991))).toEqual(53);
  expect(getMinNumBits(numToUINT64(9007199254740992))).toEqual(54);
  expect(getMinNumBits(numToUINT64(18014398509481984))).toEqual(55);
  // TODO: Fix me, maybe? Things get weird somewhere after 2^53
  // expect(getMinNumBits(numToUINT64(9223372036854775807))).toEqual(63);
  // expect(getMinNumBits(numToUINT64(9223372036854775808))).toEqual(64);
  // expect(getMinNumBits(numToUINT64(18446744073709551615))).toEqual(64);
});

test('getNumSetBits', () => {
  // Testing output against: https://golang.org/pkg/math/bits/#OnesCount64
  const basicRubric: number[] = [
    0,
    1,
    1,
    2,
    1,
    2,
    2,
    3,
    1,
    2,
    2,
    3,
    2,
    3,
    3,
    4,
    1,
    2,
    2,
    3,
    2,
    3,
    3,
    4,
    2,
    3,
    3,
    4,
    3,
    4,
    4,
    5,
    1,
    2,
    2,
    3,
    2,
    3,
    3,
    4,
    2,
    3,
    3,
    4,
    3,
    4,
    4,
    5,
    2,
    3,
    3,
    4,
    3,
    4,
    4,
    5,
    3,
    4,
    4,
    5,
    4,
    5,
    5,
    6,
    1,
    2,
    2,
    3,
    2,
    3,
    3,
    4,
    2,
    3,
    3,
    4,
    3,
    4,
    4,
    5,
    2,
    3,
    3,
    4,
    3,
    4,
    4,
    5,
    3,
    4,
    4,
    5,
    4,
    5,
    5,
    6,
    2,
    3,
    3,
    4
  ];
  const computed: number[] = [];
  for (let i = 0; i < 100; i++) {
    computed.push(getNumSetBits(numToUINT64(i)));
  }
  expect(computed).toEqual(basicRubric);
  expect(getNumSetBits(numToUINT64(4294967295))).toEqual(32);
  expect(getNumSetBits(numToUINT64(4294967296))).toEqual(1);
  expect(getNumSetBits(numToUINT64(9007199254740991))).toEqual(53);
  expect(getNumSetBits(numToUINT64(9007199254740992))).toEqual(1);
  // TODO: Fix me, maybe? Things get weird somewhere after 2^53
});

test('decomposeInclusionProof', () => {
  expect(decomposeInclusionProof(0, 0)).toEqual({ inner: 64, border: 0 });
  expect(decomposeInclusionProof(0, 1)).toEqual({ inner: 0, border: 0 });
  expect(decomposeInclusionProof(0, 2)).toEqual({ inner: 1, border: 0 });
  expect(decomposeInclusionProof(495939, 23939586939)).toEqual({
    inner: 35,
    border: 0
  });
  expect(decomposeInclusionProof(23939586939, 23939586939)).toEqual({
    inner: 1,
    border: 19
  });
  expect(decomposeInclusionProof(23939586938, 23939586939)).toEqual({
    inner: 0,
    border: 19
  });
});
