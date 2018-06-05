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
        activePage={pageTypes.SETTINGS}
        connectionStatus={connectionTypes.DISCONNECTED}
      />
    );

    expect(wrapper).toMatchSnapshot();
  });

  describe("is in the refbox mode", () => {
    it("should match snapshot", () => {
      const wrapper = shallow(
        <BottomBar
          changeActivePage={() => {}}
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

describe("When clicked,", () => {
  const realFetch = global.fetch;
  let wrapper = null;
  beforeEach(() => {
    wrapper = shallow(
      <BottomBar
        changeActivePage={() => {}}
        activePage={pageTypes.REFBOX}
        connectionStatus={connectionTypes.CONNECTED}
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
  it("the start button should pass `start` to the 'sendToServer' function", () => {
    wrapper.find("#bottom-bar__start-button").simulate("click");
  });
  it("the stop button should pass `stop` to the 'sendToServer' function", () => {
    wrapper.find("#bottom-bar__stop-button").simulate("click");
  });
});
