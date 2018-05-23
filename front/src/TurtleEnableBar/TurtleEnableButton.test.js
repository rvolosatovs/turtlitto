import React from "react";
import { shallow } from "enzyme";
import sinon from "sinon";
import TurtleEnableButton from "./TurtleEnableButton";

describe("TurtleEnableButton", () => {
  it("should match snapshot", () => {
    const wrapper = shallow(
      <TurtleEnableButton
        enabled
        id={1}
        onDisable={() => {}}
        onEnable={() => {}}
      />
    );

    expect(wrapper).toMatchSnapshot();
  });

  describe("the user clicks on the disabled button", () => {
    it("should call the `onEnable` function", () => {
      const onEnableSpy = sinon.spy();
      const onDisableSpy = sinon.spy();
      const wrapper = shallow(
        <TurtleEnableButton
          onEnable={onEnableSpy}
          onDisable={onDisableSpy}
          enabled={false}
          id={1}
        />
      );

      wrapper.simulate("click");
      expect(onEnableSpy.calledOnce).toBe(true);
      expect(onDisableSpy.notCalled).toBe(true);
    });
  });

  describe("the user clicks on the enabled button", () => {
    it("should call the `onDisable` function", () => {
      const onEnableSpy = sinon.spy();
      const onDisableSpy = sinon.spy();
      const wrapper = shallow(
        <TurtleEnableButton
          onEnable={onEnableSpy}
          onDisable={onDisableSpy}
          enabled
          id={1}
        />
      );

      wrapper.simulate("click");
      expect(onEnableSpy.notCalled).toBe(true);
      expect(onDisableSpy.calledOnce).toBe(true);
    });
  });
});
