function setupXTerm() {
  const container = document.querySelector("#terminal");

  const terminal = new Terminal();
  terminal.open(container);

  const webglAddon = new WebglAddon.WebglAddon();
  terminal.loadAddon(webglAddon);
  webglAddon.onContextLoss(() => webglAddon.dispose());

  const fitAddon = new FitAddon.FitAddon();
  terminal.loadAddon(fitAddon);
  fitAddon.fit();
  window.addEventListener("resize", () => fitAddon.fit());

  return { terminal, fitAddon, webglAddon };
}

function setupEventSource() {
  return new EventSource("/api/sse");
}

async function main() {
  const { terminal } = setupXTerm();

  const eventSource = setupEventSource();
  eventSource.addEventListener("error", (err) => {
    console.error(err);
  });
  eventSource.addEventListener("log", (event) => {
    terminal.writeln(event.data);
  });
}

main().catch((err) => console.error(err));
