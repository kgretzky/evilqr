var sessions = new Map();
var running = false;
var lastSrc = "";

const API_TOKEN = "269884a8-db69-4dcd-a47f-003b2498a72e";
const API_URL = "http://127.0.0.1:35000";

const QRCODE_ID = "846fdb44-939e-4218-aa44-7ed35dd543ce";

var QRRules = new Map();

QRRules.set("discord.com", {
    imgSelector: "[class^='qrCode-'] > img",
    buttonSelector: "",
    authSelector: "",
});
QRRules.set("www.tiktok.com", {
    imgSelector: "div[data-e2e=qr-code] > canvas",
    buttonSelector: "div[data-e2e=qr-code] > div > button",
    authSelector: "",
});
QRRules.set("accounts.binance.com", {
    imgSelector: "#wrap_app canvas",
    buttonSelector: "#wrap_app button:nth-child(2)",
    authSelector: "",
});
QRRules.set("web.telegram.org", {
    imgSelector: ".qr-canvas",
    buttonSelector: "",
    authSelector: "#page-chats"
});

class Session {
    constructor(tab) {
        this.tab = tab;
        this.running = false;
        this.imgSrc = "";
        this.imgSelector = "";
        this.buttonSelector = "";
    };

    worker() {
        var session = this;

        if (!session.running) {
            session.running = true;
        }

        chrome.tabs.sendMessage(session.tab.id, { message: "get-location" })
            .then((response) => {
                if (response !== undefined) {
                    var o = QRRules.get(response.location.host);
                    if (o !== undefined) {
                        session.imgSelector = o.imgSelector;
                        session.buttonSelector = o.buttonSelector;
                    }
                }
            });

        var foundImage = false;
        if (session.imgSelector !== "") {
            chrome.tabs.sendMessage(session.tab.id, { message: "get-image", selector: session.imgSelector })
                .then((response) => {
                    if (response !== undefined) {
                        if (response.imgSrc != "" && session.imgSrc != response.imgSrc) {
                            // new image
                            console.log("tab:" + session.tab.id + " img:" + response.imgSrc);
                            session.imgSrc = response.imgSrc;
                            foundImage = true;

                            fetch(API_URL + "/qrcode/" + QRCODE_ID, {
                                method: "PUT",
                                headers: {
                                    "Content-Type": "application/json",
                                    "Authorization": "Bearer " + API_TOKEN
                                },
                                body: JSON.stringify({ id: QRCODE_ID, source: session.imgSrc })
                            })
                                .then((response) => response.json())
                                .then((result) => {
                                    console.log("api: success:", result);
                                })
                                .catch((error) => {
                                    console.error("api: error:", error);
                                });
                        }
                    }
                })
                .catch(() => { });

            if (session.buttonSelector != "") {
                // click a reload button if available
                chrome.tabs.sendMessage(session.tab.id, { message: "click-button", selector: session.buttonSelector })
            }
        }
        if (!foundImage) {
            // check if we are not already authenticated
            if (session.authSelector !== "") {
                chrome.tabs.sendMessage(session.tab.id, { message: "item-exists", selector: session.authSelector })
                .then((response) => {
                    if (response !== undefined) {
                        if (response.exists === true) {
                            // authorized
                            session.running = false;
                            console.log("tab:" + session.tab.id + " authorized - aborting");
                            return;
                        }
                    }
                })
                .catch(() => { });
            }
        }

        function tab_callback() {
            if (chrome.runtime.lastError) {
                sessions.delete(session.tab.id);
            } else {
                setTimeout(function () { session.worker() }, 500);
            }
        }
        chrome.tabs.get(session.tab.id, tab_callback);
    }
}

chrome.runtime.onMessage.addListener(
    function (request, sender, sendResponse) {
        console.log(request);
        console.log(sender.tabId);
    }
);

function extractQR(tabId, selector) {
    console.log("tabId: " + tabId + " selector:" + selector);
    var img = document.querySelector(selector);
    if (img !== null) {
        if (img.src != lastSrc) {
            console.log(img.src);
            chrome.runtime.sendMessage({ tabId: tabId, imgSrc: img.src });
            // TODO: upload QR code

        }
    }
    return "";
}

chrome.action.onClicked.addListener((tab) => {
    if (!tab.url.includes("chrome://")) {
        var session = sessions.get(tab.id);
        if (session == undefined) {
            session = new Session(tab);
            session.imgSelector = "";
            session.buttonSelector = "";
            sessions.set(tab.id, session);
        }
        if (session.running) {
            return;
        }

        session.worker();
    }
});
