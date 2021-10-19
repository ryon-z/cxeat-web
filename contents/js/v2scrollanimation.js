let ver03 = gsap.timeline({
  scrollTrigger: {
    trigger: ".cueat_ani_title",
    start: "50% 90%",
    end: "50% 70%",
    scrub: 1,
    delay: 1,
    markers: false,
    once: false
  }
});

ver03.addLabel("start")
  .fromTo(".cueat_ani_title", { scale: 1.6, opacity: 0 }, { opacity: 1, scale: 1 })
  .addLabel("end");