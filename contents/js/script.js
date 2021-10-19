//네비 
$(document).ready(function(){
    $(".popupVideo a").click(function() {
        $(".video-popup").addClass("reveal");
        $(".video-popup .video-wrapper").remove();
        $(".video-popup").append("<div class='video-wrapper'><iframe width='560' height='315' src='https://youtube.com/embed/" + $(this).data("video") + "?rel=0&playsinline=1&autoplay=1' allow='autoplay; encrypted-media' allowfullscreen></iframe></div>")
      });
      $(".video-popup-closer").click(function() {
        $(".video-popup .video-wrapper").remove();
        $(".video-popup").removeClass("reveal");
      });
});



function fnMove(seq){
    var offset = $("#div" + seq).offset();
    $('html, body').animate({scrollTop : offset.top}, 400);
};


// 탭
$(document).ready(function(){

$('ul.tabs li').click(function(){
    var tab_id = $(this).attr('data-tab');

    $('ul.tabs li').removeClass('active');
    $('.work__projects').removeClass('active');

    $(this).addClass('active');
    $("#"+tab_id).addClass('active');
});

});
// 모달창 띄우기  
$(document).ready(function() {
    $('.project').click(function(){
        $('.modal').addClass('open');
        $("body").css("overflow", "hidden");
        $("nav").css({"opacity":"0"});
    });
   $('.modal_close').click(function(){
      $('.modal').removeClass('open');
      $("body").css("overflow", "auto");
      $("nav").css({"opacity":"1"}); 
      youtube();
    });
$('.btn_nav').click(function(){
    $('.btn_nav').removeClass('active');
    $(this).addClass('active');
    });
    
});


//gnb
// $(function(){
//     $('.gnb > li >a').on('mouseenter focus',function(){
//         var gnbindex = $('.gnb > li >a').index($(this));
//         $('.gnb .inner').removeClass('on');
//        $('.gnb .inner:eq('+ gnbindex +')').addClass('on');
//     });
//     $('header').on('mouseleave', function(){
//         $('.gnb .inner').removeClass('on');
//     } )
// });

//fixheader



var scrollFix = 0;
scrollFix = $(document).scrollTop();
fixHeader();



//윈도우창 조절시 이벤트
$(window).on('scroll resize', function(){
    scrollFix = $(document).scrollTop();
    fixHeader();
});



//고정헤더함수=> fixheader();
function fixHeader(){
    if(scrollFix > 150){
        $('header').addClass('on');
    }else{
        $('header').removeClass('on');
    }
}

//글자애니메이션 splitting 데모사이트 그대로 작성 따라하기
$(function(){
    Splitting();
});

//.top-visual 이미지 슬라이드
$(function(){
    $('.visual .slide').slick({
        arrow:true,//화살표
        dots:true,//인디케이터
        autoplay:true,//자동재생
        fade:true,//페이드인 효과
        autoplaySpeed:7000,//재생시간
        pauseOnHover:false,//호버시 멈춤 해제
        pauseOnFocus:false 
        
    });
    //두번째 슬라이드
        $('.slide2').slick({
        arrow:false,//화살표
        dots:true,//인디케이터
        autoplay:true,//자동재생
        infinity:true,
        slidesToShow: 2,
        slidesToScroll: 1,
        autoplaySpeed:6000,//재생시간
        pauseOnHover:true,//호버시 멈추도록
        pauseOnFocus:true 
        
    });
    $('.slide2 #slick-slide-control10').text("따뜻한 밥에 찌개가 먹고 싶을때");
    $('.slide2 #slick-slide-control11').text("최고의 조합 치킨과 맥주가 그리울 때");
    $('.slide2 #slick-slide-control12').text("야식대장 족발 보쌈이 생각날 때");
    $('.slide2 #slick-slide-control13').text("외식의 꽃 짜장면이 땡길 때");
});

$(function(){
    $('.animate').scrolla({
    // default
    mobile: true, // disable animation on mobiles
    once: false, // only once animation play on scroll
    animateCssVersion: 4 // used animate.css version (3 or 4)
    });
}); 

