import React from "react";
import RefboxField from ".";
import { shallow } from "enzyme";
import { mountWithTheme } from "../testUtils";

//Test_items: RefboxField.js
//Input_spec: -
//Output_spec: -
//Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
describe("RefboxField", () => {
  it("shows a cyan and magenta refbox field with all buttons", () => {
    const wrapper = mountWithTheme(<RefboxField isPenalty={false} />);
    expect(wrapper).toMatchSnapshot();
  });
  it("shows a refbox with penalty buttons", () => {
    const wrapper = mountWithTheme(<RefboxField isPenalty={true} />);
    expect(wrapper).toMatchSnapshot();
  });
});

//Test_items: RefboxField.js functionality.
//Input_spec: -
//Output_spec: -
//Envir_needs: -
describe("When clicked, a cyan", () => {
  const realFetch = global.fetch;
  let wrapper = null;
  beforeEach(() => {
    const l = window.location;
    wrapper = shallow(<RefboxField isPenalty={false} />);
    global.fetch = jest.fn().mockImplementation((url, params) => {
      expect(url).toBe(`${l.protocol}//${l.host}/api/v1/command`);
      expect(params).toMatchSnapshot();
      return Promise.resolve({ ok: true });
    });
  });

  afterEach(() => {
    global.fetch = realFetch;
  });

  it("KO button should pass `kick_off_cyan` to the 'sendToServer' function", () => {
    wrapper.find("#KO_cyan").simulate("click");
  });
  it("FK button should pass `free_kick_cyan` to the 'sendToServer' function", () => {
    wrapper.find("#FK_cyan").simulate("click");
  });
  it("GK button should pass `goal_kick_cyan` to the 'sendToServer' function", () => {
    wrapper.find("#GK_cyan").simulate("click");
  });
  it("TI button should pass `throw_in_cyan` to the 'sendToServer' function", () => {
    wrapper.find("#TI_cyan").simulate("click");
  });
  it("C button should pass `corner_cyan` to the 'sendToServer' function", () => {
    wrapper.find("#C_cyan").simulate("click");
  });
  it("P button should pass `penalty_cyan` to the 'sendToServer' function", () => {
    wrapper.find("#P_cyan").simulate("click");
  });
});

//Test_items: RefboxField index.js
//Input_spec: -
//Output_spec: -
//Envir_needs: -
describe("When clicked, a magenta", () => {
  const realFetch = global.fetch;
  let wrapper = null;
  beforeEach(() => {
    const l = window.location;
    wrapper = shallow(<RefboxField onPenalty={false} />);
    global.fetch = jest.fn().mockImplementation((url, params) => {
      expect(url).toBe(`${l.protocol}//${l.host}/api/v1/command`);
      expect(params).toMatchSnapshot();
      return Promise.resolve({ ok: true });
    });
  });

  afterEach(() => {
    global.fetch = realFetch;
  });

  it("KO button should pass `kick_off_magenta` to the 'sendToServer' function", () => {
    wrapper.find("#KO_magenta").simulate("click");
  });
  it("FK button should pass `free_kick_magenta` to the 'sendToServer' function", () => {
    wrapper.find("#FK_magenta").simulate("click");
  });
  it("GK button should pass `goal_kick_magenta` to the 'sendToServer' function", () => {
    wrapper.find("#GK_magenta").simulate("click");
  });
  it("TI button should pass `throw_in_magenta` to the 'sendToServer' function", () => {
    wrapper.find("#TI_magenta").simulate("click");
  });
  it("C button should pass `corner_magenta` to the 'sendToServer' function", () => {
    wrapper.find("#C_magenta").simulate("click");
  });
  it("P button should pass `penalty_magenta` to the 'sendToServer' function", () => {
    wrapper.find("#P_magenta").simulate("click");
  });
});
