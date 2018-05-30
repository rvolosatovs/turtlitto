import React from "react";
import RefboxField from ".";
import { shallow } from "enzyme";

describe("RefboxField", () => {
  it("shows a cyan refbox field with all buttons", () => {
    const wrapper = shallow(
      <RefboxField
        teamColor="cyan"
        onClick={(tag, color) => {}}
        isPenalty={false}
      />
    );
    expect(wrapper).toMatchSnapshot();
  });
  it("shows a magenta refbox field with all buttons", () => {
    const wrapper = shallow(
      <RefboxField
        teamColor="magenta"
        onClick={(tag, color) => {}}
        isPenalty={false}
      />
    );
    expect(wrapper).toMatchSnapshot();
  });
  it("shows a magenta refbox with penalty buttons", () => {
    const wrapper = shallow(
      <RefboxField
        teamColor="magenta"
        onClick={(tag, color) => {}}
        isPenalty={true}
      />
    );
    expect(wrapper).toMatchSnapshot();
  });
});
