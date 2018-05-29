import React from "react";
import RefboxButton from "./RefboxButton";
import renderer from "react-test-renderer";

describe("RefboxButton", () => {
  it("shows a cyan button with KO on it", () => {
    const component = renderer.create(
      <RefboxButton teamColor="cyan" onClick={() => {}} tag="KO" />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
  it("shows a magenta button with P on it", () => {
    const component = renderer.create(
      <RefboxButton teamColor="magenta" onClick={() => {}} tag="P" />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
