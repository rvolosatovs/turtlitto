import React from "react";
import InOutButton from "./InOutButton";
import { shallow } from "enzyme";

describe("InOutButton", () => {
  it("shows two buttons with go in or go out written on it", () => {
    const wrapper = shallow(<InOutButton onClick={() => {}} />);
    expect(wrapper).toMatchSnapshot();
  });
});
