import React from "react";
import BottomBar from ".";
import { shallow } from "enzyme";
import sinon from "sinon";

describe("BottomBar", () => {
  const refboxPage = "refbox";
  const settingsPage = "settings";

  describe("the user clicks on the start button", () => {
    it("should pass `start` to the `onSend` function", () => {
      const onSendSpy = sinon.spy();
      const wrapper = shallow(
        <BottomBar
          changeActivePage={() => {}}
          onSend={onSendSpy}
          activePage={refboxPage}
        />
      );

      wrapper.find("#bottom-bar__start-button").simulate("click");

      expect(onSendSpy.calledOnce).toBe(true);
      expect(onSendSpy.calledWithExactly("start")).toBe(true);
    });
  });

  describe("the user clicks on the stop button", () => {
    it("should pass `stop` to the `onSend` function", () => {
      const onSendSpy = sinon.spy();
      const wrapper = shallow(
        <BottomBar
          changeActivePage={() => {}}
          onSend={onSendSpy}
          activePage={refboxPage}
        />
      );

      wrapper.find("#bottom-bar__stop-button").simulate("click");

      expect(onSendSpy.calledOnce).toBe(true);
      expect(onSendSpy.calledWithExactly("stop")).toBe(true);
    });
  });

  describe("is in the refbox mode", () => {
    it("should match snapshot", () => {
      const wrapper = shallow(
        <BottomBar
          changeActivePage={() => {}}
          onSend={() => {}}
          activePage={refboxPage}
        />
      );

      expect(wrapper).toMatchSnapshot();
    });

    describe("the user clicks on the settings button", () => {
      it("should pass `settings` to the `changeActivePage` function", () => {
        const changeActivePageSpy = sinon.spy();
        const wrapper = shallow(
          <BottomBar
            changeActivePage={changeActivePageSpy}
            onSend={() => {}}
            activePage={refboxPage}
          />
        );

        wrapper.find("#bottom-bar__change-page-button").simulate("click");

        expect(changeActivePageSpy.calledOnce).toBe(true);
        expect(changeActivePageSpy.calledWithExactly("settings")).toBe(true);
      });
    });
  });

  describe("is in the settings mode", () => {
    it("should match snapshot", () => {
      const wrapper = shallow(
        <BottomBar
          changeActivePage={() => {}}
          onSend={() => {}}
          activePage={settingsPage}
        />
      );

      expect(wrapper).toMatchSnapshot();
    });

    describe("the user clicks on the refbox button", () => {
      it("should pass `refbox` to the `changeActivePage` function", () => {
        const changeActivePageSpy = sinon.spy();
        const wrapper = shallow(
          <BottomBar
            changeActivePage={changeActivePageSpy}
            onSend={() => {}}
            activePage={settingsPage}
          />
        );

        wrapper.find("#bottom-bar__change-page-button").simulate("click");

        expect(changeActivePageSpy.calledOnce).toBe(true);
        expect(changeActivePageSpy.calledWithExactly("refbox")).toBe(true);
      });
    });
  });
});
