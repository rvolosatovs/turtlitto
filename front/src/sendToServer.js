export default (message, destination, token) => {
  const l = window.location;
  console.log(
    `send ${message} to ${l.protocol}//${l.host}/api/v1/${destination}`
  );
  return fetch(`${l.protocol}//${l.host}/api/v1/${destination}`, {
    method: "PUT",
    headers: new Headers({
      Authorization: "Basic " + btoa(`user:${token}`)
    }),
    body: JSON.stringify(message)
  });
};
