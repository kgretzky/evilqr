var lastButtonClicked = 0;

var qrdoc = undefined;
var qrwin = undefined;
window.addEventListener('load', () => {
    qrdoc = document;
    qrwin = window;
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
        switch (request.message) {
            case "get-image":
                var img = qrdoc.querySelector(request.selector);
                if (img !== null) {
                    switch (img.tagName) {
                        case "IMG":
                            sendResponse({ imgSrc: img.src });
                            break;
                        case "CANVAS":
                            sendResponse({ imgSrc: img.toDataURL() });
                            break;
                    }
                }
                break;
            case "item-exists":
                var o = qrdoc.querySelector(request.selector);
                if (o !== null) {
                    sendResponse({ exists: true });
                } else {
                    sendResponse({ exists: false });
                }
                break;
            case "click-button":
                if (Date.now() - lastButtonClicked >= 5000) {
                    var btn = qrdoc.querySelector(request.selector);
                    if (btn !== null) {
                        lastButtonClicked = Date.now();
                        btn.click();
                    }
                }
                break;
            case "get-location":
                sendResponse({ location: qrwin.location });
                break;
        }
    });
});
