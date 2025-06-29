let inputSNILS_prev

let observer = new MutationObserver(mutationRecords => {

  const inputSNILS = document.getElementsByName("Snils")[0]

  if (inputSNILS == inputSNILS_prev) {
    return
  }

  inputSNILS_prev = inputSNILS

  if (inputSNILS !== undefined && inputSNILS.value.length > 0) {
    chrome.runtime.sendMessage({ action: "autocomplete", text: inputSNILS.value }, (response) => {
      if (response && response.suggestions) {
        // Отобразить предложенные варианты автодополнения рядом с полем ввода
        console.log("Предложения:", response.suggestions);
        // Здесь вы добавите логику для создания и отображения UI
        const passNum = document.getElementsByName("IdentityDoc.Number")[0]
        passNum.value = response.suggestions.users[0].n_doc
      }
    });
  }

});

function setObserver() {
  const ulElement = document.getElementsByTagName("body")[0]

  observer.observe(ulElement, {
    childList: true, // наблюдать за непосредственными детьми
    subtree: true, // и более глубокими потомками
    characterDataOldValue: true, // передавать старое значение в колбэк
  });
}

if (document.readyState !== 'loading') {
    setObserver();
} else {
    document.addEventListener('DOMContentLoaded', function () {
        setObserver();
    });
}
