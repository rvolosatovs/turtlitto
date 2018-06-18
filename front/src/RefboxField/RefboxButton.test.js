import React from "react";
import RefboxButton from "./RefboxButton";
import { shallowWithTheme } from "../testUtils";

describe("RefboxButton", () => {
  it("shows a cyan button with KO on it", () => {
    const wrapper = shallowWithTheme(
      <RefboxButton teamColor="cyan" onClick={() => {}} tag="KO" />
    ).dive();
    expect(wrapper).toMatchSnapshot();
  });
  it("shows a magenta button with P on it", () => {
    const wrapper = shallowWithTheme(
      <RefboxButton teamColor="magenta" onClick={() => {}} tag="P" />
    ).dive();
    expect(wrapper).toMatchSnapshot();
  });
});
