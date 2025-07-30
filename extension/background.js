// background.js
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "autocomplete") {
    // Получаем сохраненный адрес сервера
    chrome.storage.sync.get('serverAddress', (data) => {
      const serverAddress = data.serverAddress || 'http://localhost:8070'; // Используем значение по умолчанию, если не установлено

      fetch(`${serverAddress}/list`, {
        method: "POST",
        body: JSON.stringify({
          query: request.text
        })
      })
        .then(response => {
          if (!response.ok) {
            throw new Error(`Ошибка HTTP: ${response.status} ${response.statusText}`);
          }
          return response.json();
        })
        .then(data => sendResponse({ suggestions: data }))
        .catch(error => {
          console.error("Ошибка при получении автодополнения:", error);
          // Отправляем сообщение об ошибке контент-скрипту, чтобы он мог сообщить пользователю
          sendResponse({ suggestions: [], error: error.message });
        });
    });
    return true; // Указывает, что sendResponse будет вызван асинхронно
  }

  if (request.action === "getaddress") {
    // Получаем сохраненный адрес сервера
    chrome.storage.sync.get('serverAddress', (data) => {
      const serverAddress = data.serverAddress || 'http://localhost:8070'; // Используем значение по умолчанию, если не установлено

      fetch(`${serverAddress}/addr`, {
        method: "POST",
        body: JSON.stringify({
          query: request.text
        })
      })
        .then(response => {
          if (!response.ok) {
            throw new Error(`Ошибка HTTP: ${response.status} ${response.statusText}`);
          }
          return response.json();
        })
        .then(data => sendResponse({ suggestions: data }))
        .catch(error => {
          console.error("Ошибка при получении автодополнения:", error);
          // Отправляем сообщение об ошибке контент-скрипту, чтобы он мог сообщить пользователю
          sendResponse({ suggestions: [], error: error.message });
        });
    });
    return true; // Указывает, что sendResponse будет вызван асинхронно
  }
});
