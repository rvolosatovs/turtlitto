import React from "react";
import InOutButton from "./InOutButton";
import { shallow } from "enzyme";
import sendToServer from "../sendToServer";
import { mountWithTheme } from "../testUtils";

describe("InOutButton", () => {
  it("shows two buttons with go in or go out written on it", () => {
    const wrapper = mountWithTheme(<InOutButton onClick={() => {}} />);
    expect(wrapper).toMatchSnapshot();
  });
});

describe("When clicked,", () => {
  const realFetch = global.fetch;
  let wrapper = null;
  beforeEach(() => {
    wrapper = shallow(
      <InOutButton
        onClick={prop => {
          sendToServer(prop, "command");
        }}
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
  it("the Go In button should pass `go_in` to the 'sendToServer' function", () => {
    wrapper.find("#InOutButton__go-in-button").simulate("click");
  });
  it("the Go Out button should pass `go_out` to the 'sendToServer' function", () => {
    wrapper.find("#InOutButton__go-out-button").simulate("click");
  });
});
