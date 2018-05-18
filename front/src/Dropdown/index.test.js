import React from "react";
import Dropdown from ".";
import renderer from "react-test-renderer";

describe("Dropdown", () => {
  it("can be changed", () => {
    let value = "test";
    const component = renderer.create(
      <Dropdown
        currentValue={value}
        values={["test", "test2"]}
        enabled
        onChange={newValue => {
          value = newValue;
        }}
      />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();

    // Check if we can trigger onChange
    expect(value).toBe("test");
    tree.props.onChange({ target: { value: "test2" } });
    expect(value).toBe("test2");
  });
});
