let ratePlanSwiperResolution = window.innerWidth <= 768 ? 1 : 2;

// 잇템 스와이프 //
var itemSwiper = new Swiper(".item-swiper", {
  slidesPerView: 1,
  spaceBetween: 20,
  slidesPerGroup: 1,
  loopFillGroupWithBlank: false,
  loop: true,
  navigation: {
    nextEl: ".item-swiper-next",
    prevEl: ".item-swiper-prev",
  },
  breakpoints: {
    // when window width is >= 320px
    768: {
      slidesPerView: 2,
      slidesPerGroup: 2,
    },
    1024: {
      slidesPerView: 3,
      slidesPerGroup: 3,
    },
    1320: {
      slidesPerView: 4,
      slidesPerGroup: 4,
    },
  },
});

// 요금제 스와이프 //
var ratePlanSwiper = new Swiper(".ratePlan-swiper", {
  slidesPerView: 1,
  spaceBetween: 20,
  navigation: {
    nextEl: ".ratePlan-swiper-next",
    prevEl: ".ratePlan-swiper-prev",
  },
  breakpoints: {
    // when window width is >= 320px

    768: {
      slidesPerView: 1,

    },
    1024: {
      slidesPerView: 2,
    },
    1320: {
      slidesPerView: 3,
      centeredSlides: true,
      centeredSlidesBounds: true,
    },
  },
});

function handleRatePlanActive(id) {
  const ratePlanBtnList = document.querySelectorAll(".ratePlanSelect-slide");
  ratePlanSwiperResolution = window.innerWidth <= 1024 ? 1 : 2;
  for (let i = 0; i < ratePlanBtnList.length; i++) {
    const ratePlanBtn = ratePlanBtnList[i].classList;
    if (ratePlanBtn.contains("ratePlanSelect-active")) {
      ratePlanBtn.remove("ratePlanSelect-active");
    }
  }



  const targetBtnList = document.querySelectorAll(`span[id^="ratePlanBtn"`);
  for (let i = 0; i < targetBtnList.length; i++) {
    const targetBtn = targetBtnList[i];
    targetBtn.classList.remove('ratePlanSelect-active');
    
  }
  const targetBtn = document.querySelector(`#ratePlanBtn${id}`);
  targetBtn.parentElement.classList.add("ratePlanSelect-active");

  const ratePlanLabelList = document.querySelectorAll(".ratePlan-label");
  for (let i = 0; i < ratePlanLabelList.length; i++) {
    const ratePlanLabel = ratePlanLabelList[i].classList;
    if (ratePlanLabel.contains("label-active")) {
      ratePlanLabel.remove("label-active");
    }
  }
  let targetLabel = document.querySelectorAll(`#ratePlan-label${id}`);
  for (let i = 0; i < targetLabel.length; i++) {
    const target = targetLabel[i];
    target.classList.add("label-active");
  }
}

const ratePlanBtnGroup = document.querySelectorAll('span[id^="ratePlanBtn"]');

for (let i = 0; i < ratePlanBtnGroup.length; i++) {
  const ratePlanBtn = ratePlanBtnGroup[i];
  ratePlanBtn.addEventListener("click", function (e) {
    let btnIdx = e.target.id.replace(/[^0-9]/g, "");
    const resolution = window.innerWidth <= 768 ? 1: 2;
    ratePlanSwiper.slideTo(btnIdx-resolution, 500, false);
    handleRatePlanActive(btnIdx);
  });
}

const ratePlanNextBtn = document.querySelector(".ratePlan-swiper-next");
const ratePlanPrevtBtn = document.querySelector(".ratePlan-swiper-prev");

ratePlanNextBtn.addEventListener("click", function () {
  handleRatePlanActive(ratePlanSwiper.realIndex + 1);
});
ratePlanPrevtBtn.addEventListener("click", function () {
  handleRatePlanActive(ratePlanSwiper.realIndex + 1);
});
ratePlanSwiper.on("slideChange", function () {
  handleRatePlanActive(ratePlanSwiper.realIndex + 1);
});
// 이용절차 슬라이드 //
var usageSwiper = new Swiper(".usage-swiper", {
  direction: "vertical",
  slidesPerView: 1,
  spaceBetween: 800,
  effect: "fade",
  autoplay: {
    delay: 2000,
  },
  speed: 700,
});

const usageBtnList = document.querySelectorAll("button[id^='usageBtn']");

for (let i = 0; i < usageBtnList.length; i++) {
  const usageBtn = usageBtnList[i];
  usageBtn.addEventListener("click", function (e) {
    handleUsage(e.target.id);
  });
}

function handleUsageActive(id) {
  const usageNumList = document.querySelectorAll(".cueat-usage-num");
  for (let i = 0; i < usageNumList.length; i++) {
    const usageNum = usageNumList[i].classList;
    if (usageNum.contains("usage-active")) {
      usageNum.remove("usage-active");
    }
  }
  const targetBtn = document.querySelector(`#${id}`);
  targetBtn.classList.add("usage-active");
  usageSwiper.slideTo(id.replace(/[^0-9]/g, "") - 1, 500, false);
}

function handleUsageImgActive(id) {
  const usageNumList = document.querySelectorAll(".usageImg-active");
  for (let i = 0; i < usageNumList.length; i++) {
    const usageNum = usageNumList[i].classList;
    if (usageNum.contains("usageImg-active")) {
      usageNum.remove("usageImg-active");
    }
  }
  targetImg = document.querySelector(`#usageImg${id.replace(/[^0-9]/g, "")}`);
  targetImg.classList.add("usageImg-active");
}

usageSwiper.on("slideChange", function () {
  handleUsageActive(`usageNum${usageSwiper.realIndex + 1}`);
});

const usageList = document.querySelectorAll(".cueat-usage");
const usageNumList = document.querySelectorAll(".cueat-usage-num");
for (let n = 0; n < usageList.length; n++) {
  const usage = usageList[n];
  usage.addEventListener("click", function (e) {
    const usageTargetId = e.target.children[0] ? e.target.children[0].id : null;
    if (usageTargetId) {
      handleUsageActive(usageTargetId);
      handleUsageImgActive(usageTargetId);
    }
  });
}
for (let n = 0; n < usageNumList.length; n++) {
  const usageNum = usageNumList[n];
  usageNum.addEventListener("click", function (e) {
    const usageTargetId = e.target.id ? e.target.id : null;
    if (usageTargetId) {
      handleUsageActive(usageTargetId);
      handleUsageImgActive(usageTargetId);
    }
  });
}
////
//이용절차 모바일 버전//

var mobileUsageSwiper = new Swiper(".mobile-usage-swiper", {
  slidesPerView: 1,
  spaceBetween: 20,
  navigation: {
    nextEl: ".usage-swiper-next",
    prevEl: ".usage-swiper-prev",
  },
  autoplay: {
    delay: 2000,
  },
});

mobileUsageSwiper.on("slideChange", function () {
  handleMobileUsage(`m_usageBtn${mobileUsageSwiper.realIndex + 1}`);
});

function handleMobileUsage(id) {
  const mobileUsageBtnList = document.querySelectorAll("span[id^=m_usageBtn]");
  for (let n = 0; n < mobileUsageBtnList.length; n++) {
    const mobileUsageBtn = mobileUsageBtnList[n];
    mobileUsageBtn.parentNode.classList.remove("usageSelect-active");
  }
  const targetBtn = document.querySelector(`#${id}`).parentNode;
  targetBtn.classList.add("usageSelect-active");
  mobileUsageSwiper.slideTo(id.replace(/[^0-9]/g, "") - 1, 500, false);
}
const mobileUsageBtnList = document.querySelectorAll("span[id^='m_usageBtn']");

for (let i = 0; i < mobileUsageBtnList.length; i++) {
  const mobileUsageBtn = mobileUsageBtnList[i];

  mobileUsageBtn.addEventListener("click", function (e) {
    handleMobileUsage(e.target.id);
  });
}

const mobileUsageNext = document.querySelector(".usage-swiper-next");
const mobileUsagePrev = document.querySelector(".usage-swiper-prev");
mobileUsageNext.addEventListener("click", function () {
  handleMobileUsage(`m_usageBtn${mobileUsageSwiper.realIndex + 1}`);
});
mobileUsagePrev.addEventListener("click", function () {
  handleMobileUsage(`m_usageBtn${mobileUsageSwiper.realIndex + 1}`);
});

// 후기 스와이프 ///
var reviewSwiper = new Swiper(".review-swiper", {
  slidesPerView: 1,
  spaceBetween: 20,
  loopFillGroupWithBlank: true,
  loop: true,
  navigation: {
    nextEl: ".review-swiper-next",
    prevEl: ".review-swiper-prev",
  },
  breakpoints: {
    // 360:{
    //   slidesPerView: 1.1,
    //   centeredSlides: true,
    //   spaceBetween:6,
    // },
    768:{
      slidesPerView: 1,
    },
    1024: {
      slidesPerView: 2,
    },
    1320: {
      slidesPerView: 3,
    },
  },
});

// FAQ 아코디언 //

/*메뉴 선택*/

const FAQmenuList = document.querySelectorAll(".cueat-FAQ-menu");
function handleFAQmenu(id, FAQmenuList) {
  const idx = id.replace(/[^0-9]/g, "");
  for (let i = 0; i < FAQmenuList.length; i++) {
    const FAQmenu = FAQmenuList[i];
    if (FAQmenu.classList.contains("FAQmenu-active")) {
      FAQmenu.classList.remove("FAQmenu-active");
    }
  }
  const targetMenu = document.querySelector(`#FAQ-menu${idx}`);
  targetMenu.classList.add("FAQmenu-active");
}

function handleFAQList(id, FAQList) {
  const idx = id.replace(/[^0-9]/g, "");
  for (let i = 0; i < FAQList.length; i++) {
    const FAQ = FAQList[i];
    FAQ.classList.remove("FAQ-list-active");
  }
  const targetFAQ = document.querySelector(`#cueat-FAQ-list${idx}`);
  targetFAQ.classList.add("FAQ-list-active");
  // handleFAQ();
}

for (let n = 0; n < FAQmenuList.length; n++) {
  const FAQmenuList = document.querySelectorAll(".cueat-FAQ-menu");
  const FAQList = document.querySelectorAll(".cueat-FAQ-list");

  const FAQmenu = FAQmenuList[n];
  FAQmenu.addEventListener("click", function (e) {
    const targetId = e.target.id;
    handleFAQmenu(targetId, FAQmenuList);
    handleFAQList(targetId, FAQList);
  });
}
/*아코디언*/
function handleFAQ() {
  const FAQheaderList = document.querySelectorAll(".FAQ-header");
  for (let i = 0; i < FAQheaderList.length; i++) {
    const FAQheader = FAQheaderList[i];
    FAQheader.addEventListener("click", function (e) {
      const FAQheader = document.querySelector(
        `#FAQ-header${e.target.id.replace(/[^0-9]/g, "")}`
      );
      if (FAQheader.classList.contains("FAQ-active")) {
        FAQheader.classList.remove("FAQ-active");
        FAQheader.children[2].src =
          "/contents/newMain_image/cueat_plus.png";
      } else {
        FAQheader.classList.add("FAQ-active");
        FAQheader.children[2].src =
          "/contents/newMain_image/cueat_minus.png";
      }
    });
  }
}

////navBar scroll control
function handleScrollEvent() {
  const scrollY = window.scrollY;
  //nav color change
  // if(scrollY!==0){
  // 	document.querySelector('.cueat-header-nav').style.backgroundColor="#fefefe";
  // }else{
  // 	document.querySelector('.cueat-header-nav').style.backgroundColor="transparent";
  // }

  //main_count control
  if (scrollY >= 3600) {
    const cueatCount = document.querySelector(".cueat-count").textContent;
    if (cueatCount == 0) {
      main_count();
    }
  } else if (scrollY < 3600) {
    const cueatCount = document.querySelector(".cueat-count");
    cueatCount.textContent = 0;
  }
}

window.addEventListener("scroll", function () {
  handleScrollEvent();
});

document.addEventListener("DOMContentLoaded", function () {
  handleScrollEvent();
});

// main_counting

function main_count() {
  const cueatCount = document.querySelector(".cueat-count");
  var count = 0;
  var thisScriptElem = document.querySelector("#main-func")
  var countEnd = (thisScriptElem !== null) ? thisScriptElem.getAttribute("userCount") : 4000
  let counting = setInterval(function () {
    if (count > countEnd) {
      clearInterval(counting);
      cueatCount.innerHTML = new Intl.NumberFormat().format(countEnd);
      return false;
    }
    count += 212;
    cueatCount.innerHTML = new Intl.NumberFormat().format(count);
  }, 10);
}

function modalExpire(dateTime)
{
  const date = dateTime.split(' ')[0];
  const time = dateTime.split(' ')[1];
  const year = date.split('-')[0];
  const month = date.split('-')[1];
  const day = date.split('-')[2];
  const hour = time.split(':')[0];
  const min = time.split(':')[1];
  const expireDate = new Date();
  expireDate.setMinutes(min);
  expireDate.setHours(hour);
  expireDate.setDate(day);
  expireDate.setMonth(month-1);
  expireDate.setFullYear(year);
  return expireDate;
}

function preventMomentumScroll(e) {
  const body =document.querySelector("body")
  // body.classList.add("lock-screen");
}
function eventSectionControl(){
  const section = document.querySelector('.cueat-event-wrapper');
  const currentTime =new Date();
  const expireDate = modalExpire(section.dataset.timeout);
  const startDate = modalExpire(section.dataset.starttime);
  if(startDate<currentTime && currentTime < expireDate){
    section.style.display="block";
  }else{
    section.style.display="none";
  }
}
function eventModalControl(){
  const modalCloseBtn = document.querySelector('#eventModalBtn');
  const closeBtn = document.querySelector('#eventCloseBtn');
  const modal = document.querySelector('#eventModal-wrapper');
  const modalState = getCookie("eventModal");
  const currentTime =new Date();
  const modalExpireDate = modalExpire(modal.dataset.timeout);
  const modalStartDate = modalExpire(modal.dataset.starttime);
  if(modalState==null && ( modalStartDate<currentTime && currentTime < modalExpireDate)){
    window.addEventListener('scroll',preventMomentumScroll);
    window.addEventListener('touchmove',preventMomentumScroll);
    modal.style.visibility="visible";
    modal.style.opacity="1";
  }
  modalCloseBtn.addEventListener('click',function(){
    modal.style.visibility="hidden";
    modal.style.opacity="0";
    setCookie("eventModal",false, "", 10800)
    window.removeEventListener('scroll',preventMomentumScroll);
    window.removeEventListener('touchmove',preventMomentumScroll);
  })
  closeBtn.addEventListener('click',function(){
    modal.style.visibility="hidden";
    modal.style.opacity="0";
    setCookie("eventModal",false, "", 10800)
    window.removeEventListener('scroll',preventMomentumScroll);
    window.removeEventListener('touchmove',preventMomentumScroll);
  })
}


function initialize() {
  handleFAQ();
  eventModalControl();
  eventSectionControl();
}

//intialize;
initialize();
