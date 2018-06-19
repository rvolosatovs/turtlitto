export default (message, destination, session) => {
  const l = window.location;
  const msg = JSON.stringify(message);
  console.log(`send ${msg} to ${l.protocol}//${l.host}/api/v1/${destination}`);
  return fetch(`${l.protocol}//${l.host}/api/v1/${destination}`, {
    method: "POST",
    headers: new Headers({
      Authorization: "Basic " + btoa(`user:${session}`)
    }),
    body: msg
  });
};
