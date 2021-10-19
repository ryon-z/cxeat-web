function setCookie(c_name, value, domain, maxAge = 3600, path = "/") {
  var domainPhrase = domain === "" ? "" : `Domain=${domain}; `
  var c_value = `${escape(value)}; ${domainPhrase} Max-Age=${maxAge}; Path=${path}`;
  document.cookie = c_name + "=" + c_value;
}

function getCookie(cookieName) {
  var cookieValue = null;
  if (document.cookie) {
    var array = document.cookie.split(escape(cookieName) + "=");
    if (array.length >= 2) {
      var arraySub = array[1].split(";");
      cookieValue = unescape(arraySub[0]);
    }
  }
  return cookieValue;
}

function removeCookie(cookieName, domain, path="/") {
  var domainPhrase = domain === "" ? "" : `Domain=${domain}; `
  document.cookie = `${cookieName}=; ${domainPhrase}; Max-Age=0; path=${path};`;
}

function callFuncIfPressEnter(element, func) {
  element.addEventListener("keydown", function (e) {
    if (e.key === "Enter" || e.keyCode === 13) {
      func();
    }
  });
}

function customRequest(
  url,
  method,
  data,
  successFunc,
  failFunc,
  isAsync = true
) {
  var xhr = new XMLHttpRequest();

  xhr.onreadystatechange = function (e) {
    // readyStates는 XMLHttpRequest의 상태(state)를 반환
    // readyState: 4 => DONE(서버 응답 완료)
    if (xhr.readyState !== XMLHttpRequest.DONE) return;

    // status는 response 상태 코드를 반환 : 200 => 정상 응답
    if (xhr.status === 200) {
      console.log("xhr.responseText", xhr.responseText);
      successFunc(xhr.responseText);
    } else {
      console.log("Error!");
      failFunc(xhr.responseText);
    }
  };

  // 요청
  xhr.open(method, url, isAsync);
  xhr.setRequestHeader("Content-type", "application/json");
  xhr.send(JSON.stringify(data));
}

function getSelectedGoArg(goArgs, elem) {
  var splited = elem.id.split("-");
  var i = parseInt(splited[splited.length - 1], 10);
  return goArgs[i];
}

function getLabel(goArgs, codeType, codeKey) {
  for (var i = 0; i < goArgs.length; i++) {
    if (goArgs[i].CodeType === codeType && goArgs[i].CodeKey === codeKey) {
      return goArgs[i].CodeLabel;
    }
  }

  return codeKey;
}

function getItemCategoryLabel(rawCategoryLabel, goArgs, exposedSep) {
  var splited = rawCategoryLabel.split("|");
  var labels = [];

  for (var i = 0; i < splited.length; i++) {
    var label = getLabel(goArgs, "ITEM_CATEGORY", splited[i]);
    labels.push(label);
  }

  return labels.join(exposedSep);
}

function appendItemCategoryLabel(itemCategoryLabel, element) {
  var splited = itemCategoryLabel.split("|")
  for (var j=0; j<splited.length; j++) {
    var rowElem = document.createElement("div")
    rowElem.className = "product_list"
    rowElem.innerHTML = "#" + splited[j]
    element.appendChild(rowElem)
  }
}

function parseDomain() {
  var hostName = window.location.hostname
  var ipPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/

  if (hostName.slice(0, 1) === "." || hostName === "localhost" || ipPattern.test(hostName)) {
    return hostName
  }

  return "." + hostName
}

function handleLimitedInput(e, maxlength) {
  if (e.target.value.length > maxlength) {
  e.target.value = e.target.value.substr(0, maxlength);
  }
}

function isNumber(inputValue, isRequired=false) {
  var re = /[0-9]*/
  if (isRequired) {
    re = /[0-9]+/
  }

  return re.test(inputValue)
}

function checkBirthDateValid(birthDate) {
  if (!isNumber(birthDate)) {
    return false
  }
  if (birthDate.length !== 8) {
    return false
  }

  var month = parseInt(birthDate.slice(4,6), 10)
  var day = parseInt(birthDate.slice(6,8), 10)
  var monthValid = false
  var dayValid = false

  if (month >= 1 && month <= 12) {
    monthValid = true
  } else {
    return monthValid 
  }

  if ([1,3,5,7,8,10,12].includes(month) && (day >= 1 && day <= 31)) {
    dayValid = true
  } else if ([4,6,9,11].includes(month) && (day >= 1 && day <= 30)) {
    dayValid = true
  } else if (month === 2 && (day >= 1 && day <= 29)) {
    dayValid = true
  } 

  return monthValid && dayValid
}

function remove4byteEmoji(string) {
  var re = /[\u{1F004}-\u{1F9E6}]|[\u{1F600}-\u{1F9D0}]/gu
  return string.replace(re, "")
}

function applyNewLine(elem, rawString, newLineSymbol, inlineStyle=true) {
  if (elem !== null || newLineSymbol === "") {
    var splited = rawString.split(newLineSymbol)
    if (inlineStyle) {
      elem.innerHTML = splited.join("<br />")
    } else {
      for (var i=0; i<splited.length; i++) {
        var newElem = document.createElement("div")
        newElem.innerText = splited[i]
        elem.appendChild(newElem)
      }
    }
  }
}

function locateKakaoLogout(kakaoClientID) {
  var clientID = kakaoClientID
  var origin = window.location.origin
  var logoutURL = "/web/logout"
  location.href = "https://kauth.kakao.com/oauth/logout?client_id=" + clientID + "&logout_redirect_uri=" + origin + logoutURL
}
