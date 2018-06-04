import React from "react";
import { shallow } from "enzyme";
import Turtle from ".";

describe("Turtle", () => {
  it("should match snapshot", () => {
    const turtle = {
      batteryvoltage: 42,
      homegoal: "Blue home",
      id: 2,
      role: "Goalkeeper",
      teamcolor: "Cyan"
    };
    const wrapper = shallow(
      <Turtle turtle={turtle} editable={false} onChange={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
  });

  describe("the user changes dropdown values", () => {
    let wrapper = null;
    const realFetch = global.fetch;

    beforeEach(() => {
      const turtle = {
        batteryvoltage: 42,
        homegoal: "Blue home",
        id: 2,
        role: "Goalkeeper",
        teamcolor: "Cyan"
      };
      wrapper = shallow(<Turtle turtle={turtle} editable />);
      global.fetch = jest.fn().mockImplementation((url, params) => {
        expect(url).toBe("/api/v1/turtles/");
        expect(params).toMatchSnapshot();
        return Promise.resolve({ ok: true });
      });
    });

    afterEach(() => {
      global.fetch = realFetch;
    });

    it("should change the role", () => {
      wrapper.find("#turtle2__role").simulate("change", "INACTIVE");
    });

    it("should change the home", () => {
      wrapper.find("#turtle2__home").simulate("change", "Yellow home");
    });

    it("should change the team", () => {
      wrapper.find("#turtle2__team").simulate("change", "Magenta");
    });
  });
});
