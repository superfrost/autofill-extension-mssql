document.addEventListener('input', (event) => {
  // Пример: обрабатываем только текстовые поля
  if (event.target.tagName === 'INPUT' && event.target.type === 'text') {
    const inputText = event.target.value;
    if (inputText.length > 2) { // Пример: начинаем автодополнение после 2 символов
      chrome.runtime.sendMessage({ action: "autocomplete", text: inputText }, (response) => {
        if (response && response.suggestions) {
          // Отобразить предложенные варианты автодополнения рядом с полем ввода
          console.log("Предложения:", response.suggestions);
          // Здесь вы добавите логику для создания и отображения UI
        }
      });
    }
  }
});
