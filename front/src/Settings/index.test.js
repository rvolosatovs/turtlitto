import React from "react";
import { shallow } from "enzyme";
import Settings from ".";

it("renders without crashing", () => {
  const wrapper = shallow(
    <Settings
      onChange={() => {}}
      turtles={[
        {
          id: 1,
          enabled: true,
          batteryvoltage: 66,
          homegoal: "Yellow home",
          role: "INACTIVE",
          teamcolor: "Magenta"
        },
        {
          id: 2,
          enabled: false,
          batteryvoltage: 42,
          homegoal: "Yellow home",
          role: "INACTIVE",
          teamcolor: "Magenta"
        }
      ]}
      onTurtleEnableChange={() => {}}
    />
  );
  expect(wrapper).toMatchSnapshot();
});
