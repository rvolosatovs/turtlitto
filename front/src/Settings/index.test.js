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
          battery: 66,
          home: "Yellow home",
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          id: 2,
          enabled: false,
          battery: 42,
          home: "Yellow home",
          role: "INACTIVE",
          team: "Magenta"
        }
      ]}
      onTurtleEnableChange={() => {}}
    />
  );
  expect(wrapper).toMatchSnapshot();
});
