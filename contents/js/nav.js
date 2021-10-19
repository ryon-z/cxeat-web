'use strict';
//스크롤 반응하는 네비게이션 만들기
const navbar = document.querySelector('#navbar');
if (navbar !== null) {
  const navbarHeight = navbar.getBoundingClientRect().height;
  document.addEventListener("touchmove", makeNavbarDark)
  document.addEventListener('scroll', makeNavbarDark);
  function makeNavbarDark() {
    if (window.scrollY > navbarHeight) {
      navbar.classList.add('navbar--dark')
    }
    else {
      navbar.classList.remove('navbar--dark')
    }
  }

  $(function () {
    var sBtn = $('.navbar__menu > li');    //  ul > li 이를 sBtn으로 칭한다. (클릭이벤트는 li에 적용 된다.)
    sBtn.click(function () {   // sBtn에 속해 있는  a 찾아 클릭 하면.
      sBtn.removeClass("active");     // sBtn 속에 (active) 클래스를 삭제 한다.
      $(this).addClass("active"); // 클릭한 a에 (active)클래스를 넣는다.
    })
  });
  const toggleBtn = document.querySelector(".navbar__toggle");
  const modal = document.querySelector(".navbar__menu");
  const navHeight = document.querySelector("#navbar")

  toggleBtn.addEventListener('click', function () {
    modal.classList.toggle('active');
    toggleBtn.classList.toggle('open');
  });
}
