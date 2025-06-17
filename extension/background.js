// Фоновый скрипт для обработки сообщений
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === "API_SETTINGS_UPDATE") {
    // Рассылаем обновление настроек во все вкладки
    chrome.tabs.query({}, (tabs) => {
      for (const tab of tabs) {
        chrome.tabs.sendMessage(tab.id, message);
      }
    });
  }
});
