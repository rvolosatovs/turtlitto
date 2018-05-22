import React from "react";
import renderer from "react-test-renderer";
import { shallow } from "enzyme";
import sinon from "sinon";
import Turtle from ".";

describe("Turtle", () => {
  it("shows correctly", () => {
    const turtle = {
      battery: 42,
      home: "Blue home",
      id: 2,
      role: "Goalkeeper",
      team: "Cyan"
    };
    const component = renderer.create(
      <Turtle turtle={turtle} editable={false} onChange={() => {}} />
    );
    expect(component.toJSON()).toMatchSnapshot();
  });

  describe("can change the", () => {
    const turtle = {
      battery: 42,
      home: "Blue home",
      id: 2,
      role: "Goalkeeper",
      team: "Cyan"
    };
    let onChangeSpy = sinon.spy();
    const wrapper = shallow(
      <Turtle
        turtle={turtle}
        editable={true}
        onChange={(...args) => onChangeSpy(args)}
      />
    );

    beforeEach(() => {
      onChangeSpy = sinon.spy();
    });

    it("role", () => {
      wrapper.find("#turtle2__role").simulate("change", "INACTIVE");
      expect(onChangeSpy.calledOnce);
      expect(onChangeSpy.calledWithExactly("role", "INACTIVE"));
    });

    it("home", () => {
      wrapper.find("#turtle2__home").simulate("change", "Yellow home");
      expect(onChangeSpy.calledOnce);
      expect(onChangeSpy.calledWithExactly("home", "Yellow home"));
    });

    it("team", () => {
      wrapper.find("#turtle2__team").simulate("change", "Magenta");
      expect(onChangeSpy.calledOnce);
      expect(onChangeSpy.calledWithExactly("team", "Magenta"));
    });
  });
});
