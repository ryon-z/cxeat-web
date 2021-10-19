document.addEventListener("DOMContentLoaded", function () {
  var snbSwiper1Elem = document.querySelector(".snbSwiper")
  var navs = snbSwiper1Elem.querySelectorAll(".swiper-slide > button > span");
  var snbSwiper2Elem = document.querySelector(".snbSwiper2")
  var navs2 = snbSwiper2Elem.querySelectorAll(".swiper-slide > button > span");

  for (var i = 0; i < navs.length; i++) {
    navs[i].id = `s1-nav-${i}`;
    navs[i].addEventListener("click", (event) => syncNavAndBullet(event, s1Swiper));
  }

  for (var i = 0; i < navs2.length; i++) {
    navs2[i].id = `s2-nav-${i}`;
    navs2[i].addEventListener("click", (event) => syncNavAndBullet(event, s2Swiper));
  }

  function syncNavAndBullet(event, swiperElem) {
    var splited = event.target.id.split("-");
    var index = splited[splited.length - 1];

    if (swiperElem !== null) {
      swiperElem.slideTo(index, 500);
    }
  }
});

var s1Swiper = new Swiper(".s1", {
  slidesPerView: 1.2,
  spaceBetween: 20,
  breakpoints: {
    480: {
      slidesPerView: 3,
      spaceBetween: 21,
    },
  },
});

var s2Swiper = new Swiper(".s2", {
  slidesPerView: 1,
  spaceBetween: 20,
  breakpoints: {
    480: {
      slidesPerView: 1,
      spaceBetween: 21,
    },
  },
});

s1Swiper.on("activeIndexChange", function () {
  var navElem = document.querySelector(`#s1-nav-${s1Swiper.activeIndex}`);
  if (navElem !== null) {
    navElem.click();
  }
});
s2Swiper.on("activeIndexChange", function () {
  var navElem = document.querySelector(`#s2-nav-${s2Swiper.activeIndex}`);
  if (navElem !== null) {
    navElem.click();
  }
});

var snbSwiperClass1 = ".snbSwiper"
var swiper = new Swiper(snbSwiperClass1, {
  slidesPerView: "auto",
  preventClicks: true,
  preventClicksPropagation: false,
  observer: true,
  observeParents: true,
});
var $snbSwiperItem = $(`${snbSwiperClass1} .swiper-wrapper .swiper-slide button`);
$snbSwiperItem.click(function () {
  var target = $(this).parent();
  $snbSwiperItem.parent().removeClass("on");
  target.addClass("on");
  muCenter(target, snbSwiperClass1);
});

var snbSwiperClass2 = ".snbSwiper2"
var swiper = new Swiper(snbSwiperClass2, {
  slidesPerView: "auto",
  preventClicks: true,
  preventClicksPropagation: false,
  observer: true,
  observeParents: true,
});
var $snbSwiperItem2 = $(`${snbSwiperClass2} .swiper-wrapper .swiper-slide button`);
$snbSwiperItem2.click(function () {
  var target = $(this).parent();
  $snbSwiperItem2.parent().removeClass("on");
  target.addClass("on");
  muCenter(target, snbSwiperClass2);
});

function muCenter(target, cssSelector) {
  var snbwrap = $(`${cssSelector} .swiper-wrapper`);
  var targetPos = target.position();
  var box = $(cssSelector);
  var boxHarf = box.width() / 2;
  var pos;
  var listWidth = 0;

  snbwrap.find(".swiper-slide").each(function () {
    listWidth += $(this).outerWidth();
  });

  var selectTargetPos = targetPos.left + target.outerWidth() / 2;
  if (selectTargetPos <= boxHarf) {
    // left
    pos = 0;
  } else if (listWidth - selectTargetPos <= boxHarf) {
    //right
    pos = listWidth - box.width();
  } else {
    pos = selectTargetPos - boxHarf;
  }

  setTimeout(function () {
    snbwrap.css({
      transform: "translate3d(" + pos * -1 + "px, 0, 0)",
      "transition-duration": "500ms",
    });
  }, 200);
}
