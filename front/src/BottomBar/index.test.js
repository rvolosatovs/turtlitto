import React from "react";
import BottomBar from ".";
import { shallow } from "enzyme";
import sinon from "sinon";
import connectionTypes from "./connectionTypes";
import pageTypes from "./pageTypes";

describe("BottomBar", () => {
  it("should match snapshot when connected", () => {
    const wrapper = shallow(
      <BottomBar
        changeActivePage={() => {}}
        onSend={() => {}}
        activePage={pageTypes.SETTINGS}
        connectionStatus={connectionTypes.CONNECTED}
      />
    );

    expect(wrapper).toMatchSnapshot();
  });

  it("should match snapshot when connecting", () => {
    const wrapper = shallow(
      <BottomBar
        changeActivePage={() => {}}
        onSend={() => {}}
        activePage={pageTypes.SETTINGS}
        connectionStatus={connectionTypes.CONNECTING}
      />
    );

    expect(wrapper).toMatchSnapshot();
  });

  it("should match snapshot when disconnected", () => {
    const wrapper = shallow(
      <BottomBar
        changeActivePage={() => {}}
        onSend={() => {}}
        activePage={pageTypes.SETTINGS}
        connectionStatus={connectionTypes.DISCONNECTED}
      />
    );

    expect(wrapper).toMatchSnapshot();
  });

  describe("the user clicks on the start button", () => {
    it("should pass `start` to the `onSend` function", () => {
      const onSendSpy = sinon.spy();
      const wrapper = shallow(
        <BottomBar
          changeActivePage={() => {}}
          onSend={onSendSpy}
          activePage={pageTypes.REFBOX}
          connectionStatus={connectionTypes.CONNECTED}
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
          activePage={pageTypes.REFBOX}
          connectionStatus={connectionTypes.CONNECTED}
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
          activePage={pageTypes.REFBOX}
          connectionStatus={connectionTypes.CONNECTED}
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
            activePage={pageTypes.REFBOX}
            connectionStatus={connectionTypes.CONNECTED}
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
          activePage={pageTypes.SETTINGS}
          connectionStatus={connectionTypes.CONNECTED}
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
            activePage={pageTypes.SETTINGS}
            connectionStatus={connectionTypes.CONNECTED}
          />
        );

        wrapper.find("#bottom-bar__change-page-button").simulate("click");

        expect(changeActivePageSpy.calledOnce).toBe(true);
        expect(changeActivePageSpy.calledWithExactly("refbox")).toBe(true);
      });
    });
  });
});
