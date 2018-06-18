import React from "react";
import Battery from ".";
import { mountWithTheme } from "../testUtils";

describe("Battery", () => {
  it("shows a full battery on 99%", () => {
    const wrapper = mountWithTheme(<Battery percentage={99} />);
    expect(wrapper).toMatchSnapshot();
  });

  it("shows a 3/4 full battery on 75%", () => {
    const wrapper = mountWithTheme(<Battery percentage={75} />);
    expect(wrapper).toMatchSnapshot();
  });

  it("shows a half full battery on 50%", () => {
    const wrapper = mountWithTheme(<Battery percentage={50} />);
    expect(wrapper).toMatchSnapshot();
  });

  it("shows a 1/4 full battery on 25%", () => {
    const wrapper = mountWithTheme(<Battery percentage={25} />);
    expect(wrapper).toMatchSnapshot();
  });

  it("shows an empty battery on 0%", () => {
    const wrapper = mountWithTheme(<Battery percentage={0} />);
    expect(wrapper).toMatchSnapshot();
  });
});
