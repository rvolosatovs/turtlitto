import React from "react";

/**
 * Show a list of items, but only the items that are enabled.
 * Author: H.E. van der Laan
 *
 * Arguments:
 * - The component to be wrapped
 * - An array of props to be passed to the wrapped components.
 *
 * One WrappedComponent is created for each element of the data array for which the element "enabled" is truthy.
 *
 * Usage warning: Each element of the array data should have an unique element id which will be used as a reconciling hint for React.
 * See https://reactjs.org/warnings/special-props.html and https://reactjs.org/docs/lists-and-keys.html#keys for more details.
 */
export default (WrappedComponent, data) => {
  return (
    <div>
      {data
        .filter(item => {
          return item["enabled"];
        })
        .map(item => {
          return <WrappedComponent key={item.id} {...item} />;
        })}
    </div>
  );
};
