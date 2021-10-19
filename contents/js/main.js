'use strict';
//탭버튼 액티브
$(function () {
  var tabBtn = $('.category__btn');    //  ul > li 이를 sBtn으로 칭한다. (클릭이벤트는 li에 적용 된다.)
  tabBtn.click(function () {   // sBtn에 속해 있는  a 찾아 클릭 하면.
    tabBtn.removeClass("active");     // sBtn 속에 (active) 클래스를 삭제 한다.
    $(this).addClass("active"); // 클릭한 a에 (active)클래스를 넣는다.
  })
});

// 탭
// $(document).ready(function(){
//   $('ul.tabs li').click(function(){
//       var tab_id = $(this).attr('data-tab');

//       $('ul.tabs li').removeClass('active');
//       $('.work__projects').removeClass('active');

//       $(this).addClass('active');
//       $("#"+tab_id).addClass('active');
//   });
// });

//설문조사 탭
$(document).ready(function () {
  $('.purchase_tabs ul li').click(function () {
    var tab_id = $(this).attr('data-tab');

    $('.purchase_tabs li').removeClass('active');
    $('.work__projects').removeClass('active');

    $(this).addClass('active');
    $("#" + tab_id).addClass('active');
  });
});
//설문조사 탭
$(document).ready(function () {
  $('.qna_tabs ul li').click(function () {
    var tab_id = $(this).attr('data-tab');

    $('.qna_tabs li').removeClass('active');
    $('.work__projects').removeClass('active');

    $(this).addClass('active');
    $("#" + tab_id).addClass('active');
  });
});

//.top-visual 이미지 슬라이드
$(function () {
  if ($('.visual .slide').length) {
    $('.visual .slide').slick({
      infinite: true,
      slidesToShow: 1,
      slidesToScroll: 1,
      arrows: false,//화살표
      dots: true,//인디케이터
      autoplay: true,//자동재생
      cssEase: 'linear',
      fade: true,//페이드인 효과
      autoplaySpeed: 9000,//재생시간
      pauseOnHover: false,//호버시 멈춤 해제
      pauseOnFocus: false,
      pauseOnDotsHover: false
    });
  }
});

$(function () {
  if ($('.photo_review .slide').length) {
    $('.photo_review .slide').slick({
      infinite: true,
      slidesToShow: 3,
      arrows: true,
      autoplay: true,
      autoplaySpeed: 5000,
      waitForAnimate: false,
      pauseOnFocus: false,
      pauseOnHover: false,
      swipeToSlide: true,
      swipe: true,

      responsive: [
        {
          breakpoint: 1024,
          settings: {
            slidesToShow: 3,
            slidesToScroll: 3,
            infinite: true,

          }
        },
        {
          breakpoint: 600,
          settings: {
            slidesToShow: 2,
            slidesToScroll: 2
          }
        },
        {
          breakpoint: 480,
          settings: {
            slidesToShow: 1,
            autospeed: 6000,
            infinite: true,
            swipeToSlide: true,
            waitForAnimate: false,
            pauseOnFocus: false,
            pauseOnHover: false,
            swipe: true

          }
        }
        // You can unslick at a given breakpoint now by adding:
        // settings: "unslick"
        // instead of a settings object
      ]
    });
  }
});

$(function () {
  if ($('.cueat_desc .slide').length) {
    $('.cueat_desc .slide').slick({
      infinite: true,
      slidesToShow: 1,
      slidesToScroll: 1,
      arrows: false,//화살표
      dots: true,//인디케이터
      autoplay: true,//자동재생
      cssEase: 'linear',
      fade: true,//페이드인 효과
      autoplaySpeed: 9000,//재생시간
      pauseOnHover: false,//호버시 멈춤 해제
      pauseOnFocus: false,
      pauseOnDotsHover: false
    });
  }
});

// nav-slider 컨트롤
// 네비 슬라이드
$(function () {
  if ($('.intro_slider_nav').length) {
    $('.intro_slider_nav').slick({
      slidesToShow: 5,
      slidesToScroll: 1,
      asNavFor: '.intro_slider_for',
      dots: false,
      centerMode: true,
      focusOnSelect: true,
      infinite: false,
      arrows: false,
    })
  }
});

// 컨텐츠 슬라이드
$(function () {
  if ($('.intro_slider_for').length) {
    $('.intro_slider_for').slick({
      slidesToShow: 1.8,
      slidesToScroll: 1,
      arrows: false,
      fade: false,
      asNavFor: '.intro_slider_nav',

      // responsive: [
      //   {
      //     breakpoint: 1024,
      //     settings: {
      //       slidesToShow: 3,
      //       slidesToScroll: 3,
      //       infinite: true,
      //       dots: true
      //     }
      //   },
      //   {
      //     breakpoint: 600,
      //     settings: {
      //       slidesToShow: 2,
      //       slidesToScroll: 2
      //     }
      //   },
      //   {
      //     breakpoint: 480,
      //     settings: {
      //       infinite: true,
      //       slidesToShow: 1.2,
      //       autoplay: true,
      //       autoplaySpeed: 0,
      //       speed: 25000,
      //       cssEase: 'linear',
      //       infinite: true,
      //       arrows: false,
      //       touchMove: true,
      //       swipeToSlide: true,
      //       swipe: true
      //     }
      //   }
      // ]
    })
  }
});