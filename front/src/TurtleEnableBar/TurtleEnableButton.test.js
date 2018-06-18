import React from "react";
import TurtleEnableButton from "./TurtleEnableButton";
import { shallowWithTheme } from "../testUtils";

describe("TurtleEnableButton", () => {
  it("should match snapshot", () => {
    const wrapper = shallowWithTheme(
      <TurtleEnableButton enabled id={1} onTurtleEnableChange={() => {}} />
    ).dive();

    expect(wrapper).toMatchSnapshot();
  });
});
