export default (message, destination) => {
  const l = window.location;
  console.log(
    `send ${message} to ${l.protocol}//${l.host}/api/v1/${destination}`
  );
  return fetch(`${l.protocol}//${l.host}/api/v1/${destination}`, {
    method: "PUT",
    headers: {},
    body: JSON.stringify(message)
  });
};
