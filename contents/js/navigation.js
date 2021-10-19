
const navMenuIcon = document.querySelector("#navMenuToggle");

if (navMenuIcon !== null) {
  navMenuIcon.addEventListener("click", function (e) {
    const mobileMenu = document.querySelector(".cueat-mobileMenu-list");
    const targetBtn = e.target.classList;
    if (targetBtn.contains("menu-active")) {
      targetBtn.remove("menu-active");
      mobileMenu.style.transform = "translate(100%,0)";
      e.target.src = "/contents/newMain_image/hamburger_icon.png";
    } else {
      targetBtn.add("menu-active");
      mobileMenu.style.transform = "translate(0,0)";
      e.target.src = "/contents/newMain_image/close.png";
    }
  });
}

const linkToLoginList = document.querySelectorAll(".linkToLogin");
for (let i = 0; i < linkToLoginList.length; i++) {
  const linkToLogin = linkToLoginList[i];
  linkToLogin.addEventListener("click", function () {
    var currentPath = window.location.pathname + window.location.search;
    setCookie("pageAfterLogin", currentPath, "", 3600, "/");
    location.href = "/web/login";
  });
}
