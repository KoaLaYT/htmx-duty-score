window.addEventListener("load", () => {
  let key = "";
  const evtSource = new EventSource("http://127.0.0.1:41987/_dev/sse");
  evtSource.onerror = (error) => {
    console.error(error);
    evtSource.close();
  };
  evtSource.onmessage = (event) => {
    if (key && key !== event.data) {
      window.location.reload();
    }
    key = event.data;
  };
});
