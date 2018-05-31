import React from "react";
import { shallow } from "enzyme";
import TurtleEnableButton from "./TurtleEnableButton";

describe("TurtleEnableButton", () => {
  it("should match snapshot", () => {
    const wrapper = shallow(
      <TurtleEnableButton enabled id={1} onTurtleEnableChange={() => {}} />
    );

    expect(wrapper).toMatchSnapshot();
  });
});
