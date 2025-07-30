let inputSNILS_prev

const nivs = Object.getOwnPropertyDescriptor(
  window.HTMLInputElement.prototype,
  "value"
).set

const ntavs = Object.getOwnPropertyDescriptor(
  window.HTMLTextAreaElement.prototype,
  "value"
).set

let observer = new MutationObserver(mutationRecords => {

  const inputSNILS = document.getElementsByName("Snils")[0]

  if (inputSNILS == inputSNILS_prev) {
    return
  }

  inputSNILS_prev = inputSNILS

  if (inputSNILS !== undefined && inputSNILS.value.length > 0) {
    chrome.runtime.sendMessage({ action: "autocomplete", text: inputSNILS.value }, async (response) => {
      if (response && response.suggestions) {
        // Отобразить предложенные варианты автодополнения рядом с полем ввода
        console.log("Предложения:", response.suggestions);
        chrome.runtime.sendMessage({ action: "getaddress", text: response.adres }, (response) => {
          if (response && response.suggestions) {
            // Отобразить предложенные варианты автодополнения рядом с полем ввода
            console.log("Адрес:", response.suggestions);
            // Здесь вы добавите логику для создания и отображения UI
            const flat = document.getElementsByName("Address.Flat")[0]
            flat.setAttribute("value", response.suggestions.flat)

            const house = document.getElementsByName("Address.Building")[0]
            house.setAttribute("value", response.suggestions.house)

            const index = document.getElementsByName("Address.ZipCode")[0]
            index.setAttribute("value", response.suggestions.post_index)
            
            const street = document.getElementsByName("Address.Street")[0]
            street.setAttribute("value", response.suggestions.street)

            const city = document.getElementsByName("Address.Place")[0]
            city.setAttribute("value", response.suggestions.town_name)

            const subject = document.getElementsByName("Address.TerritorySubject")[0]
            subject.setAttribute("value", response.suggestions.rf_subj)

            const district = document.getElementsByName("Address.District")[0]
            district.setAttribute("value", response.suggestions.rf_subj_reg)
          }
        });
        // Здесь вы добавите логику для создания и отображения UI
        const docSeries = document.getElementsByName("IdentityDoc.Series")[0]
        nivs.call(docSeries, response.suggestions.users[0].s_doc)
        docSeries.dispatchEvent(new Event("change", {bubbles: true}))

        const docNum = document.getElementsByName("IdentityDoc.Number")[0]
        nivs.call(docNum, response.suggestions.users[0].n_doc)
        docNum.dispatchEvent(new Event("change", {bubbles: true}))

        const docDate = document.getElementsByName("IdentityDoc.IssueDate")[0]
        const docDate1 = new Date(response.suggestions.users[0].dateDoc)
        const day = docDate1.getDate()
        const mon = docDate1.getMonth()+1 < 10 ? "0"+(docDate1.getMonth()+1) : docDate1.getMonth()+1
        const year = docDate1.getFullYear()
        const docDate2 = day+"."+mon+"."+year
        nivs.call(docDate, docDate2)
        docDate.dispatchEvent(new Event("change", {bubbles: true}))

        const docOrg = document.getElementsByName("IdentityDoc.IssueOrgName")[0]
        ntavs.call(docOrg, response.suggestions.users[0].docIssuedBy)
        docOrg.dispatchEvent(new Event("change", {bubbles: true}))

        const phone = document.getElementsByName("Phone")[0]
        nivs.call(phone, response.suggestions.users[0].contactMPhone)
        phone.dispatchEvent(new Event("change", {bubbles: true}))

        const email = document.getElementsByName("Email")[0]
        nivs.call(email, response.suggestions.users[0].contactEmail)
        email.dispatchEvent(new Event("change", {bubbles: true}))

        const docType = document.getElementsByName("IdentityDoc.IdentityCardTypeId.Id")[0]
        const docTypeDiv = document.getElementsByClassName(" css-gdh54p-singleValue")[2]
        const doc = {}
        switch(response.suggestions.users[0].rf_TYPEDOCID) {
          case "3": doc.code = 1; doc.name = "Паспорт гражданина Российской Федерации"; break;
          case "5": doc.code = 2; doc.name = "Свидетельство о рождении"; break;
        }
        docTypeDiv.innerText = doc.name
        nivs.call(docType, doc.code)
        docType.dispatchEvent(new Event("change", {bubbles: true}))

        
        
      }
    });
  }

  const list = document.getElementsByTagName("li")
  const listArray = [...list]
  const tabs = listArray.filter(l => l.classList.contains("is-active"))[0]
  
  if (tabs.title === "Сведения о результатах предыдущей медико-социальной экспертизы") {
    const Profession = document.getElementById("EducationInfo.Profession")
    Profession.value = "+"

    const OrgName = document.getElementById("EducationInfo.OrgName")
    OrgName.value = "+"

    const Street = document.getElementsByName("EducationInfo.OrgAddress.Street")[0]
    Street.value = "+"

    const Building = document.getElementById("EducationInfo.OrgAddress.Building")
    Building.value = "+"

    const ZipCode = document.getElementsByName("EducationInfo.OrgAddress.ZipCode")[0]
    ZipCode.value = "+"

    const LevelValue = document.getElementById("EducationInfo.LevelValue")
    LevelValue.value = "+"

    const MainProfession = document.getElementById("ProfInfo.MainProfession")
    MainProfession.value = "+"

    const Qualification = document.getElementById("ProfInfo.Qualification")
    Qualification.value = "+"

    const JobExperience = document.getElementById("ProfInfo.JobExperience")
    JobExperience.value = "+"

    const Speciality = document.getElementById("ProfInfo.CurrentJob.Speciality")
    Speciality.value = "+"

    const Position = document.getElementById("ProfInfo.CurrentJob.Position")
    Position.value = "+"

    const JobPlace = document.getElementById("ProfInfo.JobPlace")
    JobPlace.value = "+"

    const LaborConditions = document.getElementById("ProfInfo.LaborConditions")
    LaborConditions.value = "+"

    const JobProfession = document.getElementById("ProfInfo.CurrentJob.Profession")
    JobProfession.value = "+"

    const JobStreet = document.getElementById("ProfInfo.JobAddress.Street")
    JobStreet.value = "+"

    const JobBuilding = document.getElementById("ProfInfo.JobAddress.Building")
    JobBuilding.value = "+"

    const JobZipCode = document.getElementsByName("ProfInfo.JobAddress.ZipCode")[0]
    JobZipCode.value = "+"
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

let body

if (document.readyState !== 'loading') {
    setObserver();
    body = document.getElementsByTagName("body")[0]
} else {
    document.addEventListener('DOMContentLoaded', function () {
        setObserver();
        body = document.getElementsByTagName("body")[0]
    });
}
