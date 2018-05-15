import React from "react";
import Dropdown from "./Dropdown";
import renderer from "react-test-renderer";

test("Dropdowns can be changed", () => {
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

test("Dropdowns can be disabled", () => {
  let value = "test";
  const component = renderer.create(
    <Dropdown
      currentValue={value}
      values={["test", "test2"]}
      enabled={false}
      onChange={newValue => {
        value = newValue;
      }}
    />
  );
  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();

  // Check if we can't trigger onChange
  expect(value).toBe("test");
  tree.props.onChange({ target: { value: "test2" } });
  expect(value).toBe("test");
});
