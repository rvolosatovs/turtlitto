import React from "react";
import { shallow } from "enzyme";
import sinon from "sinon";
import Turtle from ".";

describe("Turtle", () => {
  it("should match snapshot", () => {
    const turtle = {
      battery: 42,
      home: "Blue home",
      id: 2,
      role: "Goalkeeper",
      team: "Cyan"
    };
    const wrapper = shallow(
      <Turtle turtle={turtle} editable={false} onChange={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
  });

  describe("the user changes dropdown values", () => {
    let onChangeSpy = null;
    let wrapper = null;

    beforeEach(() => {
      const turtle = {
        battery: 42,
        home: "Blue home",
        id: 2,
        role: "Goalkeeper",
        team: "Cyan"
      };
      onChangeSpy = sinon.spy();
      wrapper = shallow(
        <Turtle
          turtle={turtle}
          editable={true}
          onChange={(...args) => onChangeSpy(args)}
        />
      );
    });

    it("should change the role", () => {
      wrapper.find("#turtle2__role").simulate("change", "INACTIVE");
      expect(onChangeSpy.calledOnce);
      expect(onChangeSpy.calledWithExactly("role", "INACTIVE"));
    });

    it("should change the home", () => {
      wrapper.find("#turtle2__home").simulate("change", "Yellow home");
      expect(onChangeSpy.calledOnce);
      expect(onChangeSpy.calledWithExactly("home", "Yellow home"));
    });

    it("should change the team", () => {
      wrapper.find("#turtle2__team").simulate("change", "Magenta");
      expect(onChangeSpy.calledOnce);
      expect(onChangeSpy.calledWithExactly("team", "Magenta"));
    });
  });
});
