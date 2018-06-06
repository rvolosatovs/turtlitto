import { css } from "styled-components";

const sizes = {
  xl: 1170,
  lg: 992,
  md: 768,
  sm: 576,
  xs: 376
};

export const media = Object.keys(sizes).reduce((accumulator, label) => {
  const size = sizes[label];
  accumulator[label] = (...args) => css`
    @media (min-width: ${size}px) {
      ${css(...args)};
    }
  `;
  return accumulator;
}, {});
