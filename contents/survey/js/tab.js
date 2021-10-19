// 배송지 탭
$(".delivery_tabs ul li").click(function () {
  var tab_id = $(this).attr("data-tab");

  $(".delivery_tabs li").removeClass("active");
  $(".work__projects").removeClass("active");

  $(this).addClass("active");
  $("#" + tab_id).addClass("active");
});

// 카드 탭
$(".card_tabs ul li").click(function () {
  var tab_id2 = $(this).attr("data-tab");

  $(".card_tabs li").removeClass("active");
  $(".work__projects2").removeClass("active");

  $(this).addClass("active");
  $("#" + tab_id2).addClass("active");
});
