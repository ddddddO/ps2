const go = new Go();
let mod, instance;
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    mod = result.module;
    instance = result.instance;

    console.clear();
    go.run(instance);
    instance = WebAssembly.instantiate(mod, go.importObject);
});

const callPS2 = () => {
    ps2Run();
};

const getOutputType = () => {
  const elms = document.getElementsByName("output_type");
  elms.forEach((elm) => {
    if (elm.checked) {
      console.debug('👺<checked! ');
      console.dir(elm, {depth: null})
    }
  })
}

const clearInput = () => {
  document.getElementById("input").value = "";
};

const clearOutput = () => {
  const elm = document.getElementById("redered_json");
  if (elm !== null) {
    elm.remove();
  }
}

const clearFn = () => {
  clearInput();
  clearOutput();
}

const copyToClipboard = () => {
  const calender = document.getElementById("redered_json");
  if (calender === null) {
    return;
  }
  const clipboard = window.navigator.clipboard;
  clipboard.writeText(calender.textContent);
};
