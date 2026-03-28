chrome.downloads.onChanged.addListener((delta) => {
    if (!delta.state || delta.state.current !== "complete") return;

    chrome.downloads.search({ id: delta.id }, (results) => {
        if (!results || results.length === 0) return;

        let item = results[0];
        let filename = item.filename;

        let savedExt = filename.split('.').pop().toLowerCase();

        let port = chrome.runtime.connectNative("ffconvert_helper");

port.onMessage.addListener((response) => {
    console.log("Helper:", response);
});

port.onDisconnect.addListener(() => {
    if (chrome.runtime.lastError) {
        console.log("Native host closed:", chrome.runtime.lastError.message);
    }
});


        port.postMessage({
            input: filename,
            target_ext: savedExt,
            force: false
        });
    });
});
