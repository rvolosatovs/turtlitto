import React from "react";
import RefboxButton from "./RefboxButton";
import { shallow } from "enzyme";

describe("RefboxButton", () => {
  it("shows a cyan button with KO on it", () => {
    const wrapper = shallow(
      <RefboxButton teamColor="cyan" onClick={() => {}} tag="KO" />
    );
    expect(wrapper).toMatchSnapshot();
  });
  it("shows a magenta button with P on it", () => {
    const wrapper = shallow(
      <RefboxButton teamColor="magenta" onClick={() => {}} tag="P" />
    );
    expect(wrapper).toMatchSnapshot();
  });
});
