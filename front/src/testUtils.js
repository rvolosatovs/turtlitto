import React from "react";
import { ThemeProvider } from "styled-components";
import theme from "./theme";
import { shallow, mount } from "enzyme";

export const shallowWithTheme = tree => {
  const context = shallow(<ThemeProvider theme={theme} />)
    .instance()
    .getChildContext();
  return shallow(tree, { context });
};

export const mountWithTheme = tree => {
  const context = shallow(<ThemeProvider theme={theme} />)
    .instance()
    .getChildContext();

  return mount(tree, {
    context,
    childContextTypes: ThemeProvider.childContextTypes
  });
};
