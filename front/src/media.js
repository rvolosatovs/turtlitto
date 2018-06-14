import { css } from "styled-components";

export const screenSizes = {
  xl: 1170,
  lg: 992,
  md: 768,
  sm: 576,
  xs: 376
};

export default Object.keys(screenSizes).reduce((accumulator, label) => {
  const size = screenSizes[label];
  accumulator[label] = (...args) => css`
    @media (min-width: ${size}px) {
      ${css(...args)};
    }
  `;
  return accumulator;
}, {});
