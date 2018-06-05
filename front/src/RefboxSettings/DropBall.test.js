import React from "react";
import DropBall from "./DropBall";
import { shallow } from "enzyme";
import sendToServer from "../sendToServer";

describe("DropBall", () => {
  it("shows a button with DB written on it", () => {
    const wrapper = shallow(<DropBall onClick={() => {}} />);
    expect(wrapper).toMatchSnapshot();
  });
});

describe("When clicked, the Drop Ball button", () => {
  const realFetch = global.fetch;
  let wrapper = null;
  beforeEach(() => {
    wrapper = shallow(
      <DropBall
        onClick={() => {
          sendToServer("dropped_ball", "command");
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
  it("should pass `dropped_ball` to the 'sendToServer' function", () => {
    wrapper.find("#drop-ball-button").simulate("click");
  });
});
