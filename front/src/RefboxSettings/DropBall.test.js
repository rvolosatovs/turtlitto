import React from "react";
import DropBall from "./DropBall";
import { shallow } from "enzyme";

describe("DropBall", () => {
  it("shows a button with DB written on it", () => {
    const wrapper = shallow(<DropBall onClick={() => {}} />);
    expect(wrapper).toMatchSnapshot();
  });
});
