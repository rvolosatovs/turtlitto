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

describe("When the role assigner dropdown is changed to", () => {
  const realFetch = global.fetch;
  let wrapper = null;
  beforeEach(() => {
    wrapper = shallow(
      <Settings
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
      />
    );
    const l = window.location;
    global.fetch = jest.fn().mockImplementation((url, params) => {
      expect(url).toBe(`${l.protocol}//${l.host}/api/v1/command`);
      expect(params).toMatchSnapshot();
      return Promise.resolve({ ok: true });
    });
  });

  afterEach(() => {
    global.fetch = realFetch;
  });
  it("'Role Assigner on', it should pass 'role_assigner_on' to sendToServer", () => {
    wrapper
      .find("#settings_role-dropdown")
      .simulate("change", "Role assigner on");
  });
  it("'Role Assigner off', it should pass 'role_assigner_off' to sendToServer", () => {
    wrapper
      .find("#settings_role-dropdown")
      .simulate("change", "Role assigner off");
  });
  it("'Pass demo', it should pass 'pass_demo' to sendToServer", () => {
    wrapper.find("#settings_role-dropdown").simulate("change", "Pass demo");
  });
  it("'Penalty mode', it should pass 'penalty_demo' to sendToServer", () => {
    wrapper.find("#settings_role-dropdown").simulate("change", "Penalty mode");
  });
  it("'Ball Handling demo', it should pass 'ball_handling_demo' to sendToServer", () => {
    wrapper
      .find("#settings_role-dropdown")
      .simulate("change", "Ball Handling demo");
  });
});
