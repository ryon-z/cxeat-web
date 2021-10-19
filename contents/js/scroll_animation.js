const boxes = gsap.utils.toArray('.en_title_p');
boxes.forEach((box, i) => {
  const anim = gsap.fromTo(box, { autoAlpha: 0, y: 50 }, { duration: 1, autoAlpha: 1, y: 0 });
  ScrollTrigger.create({
    trigger: box,
    animation: anim,
    start: "0 80%",
    end: "0 50%",
    markers: false,
    toggleActions: 'play none none none',
    once: false,
    ease: "power1.inOut"
  });
});

const koTexts = gsap.utils.toArray('.ko_text_p');
koTexts.forEach((ko, i) => {
  const anim = gsap.fromTo(ko, { autoAlpha: 0, y: 50 }, { duration: 1, autoAlpha: 1, y: 0 });
  ScrollTrigger.create({
    trigger: ko,
    animation: anim,
    start: "0 80%",
    end: "0 50%",
    markers: false,
    toggleActions: 'play none none none',
    once: false,
    ease: "power1.inOut"
  });
});



let tl = gsap.timeline({
  // yes, we can add it to an entire timeline!
  scrollTrigger: {
    trigger: ".curation_div1_img003",
    start: "0% 80%", // when the top of the trigger hits the top of the viewport
    end: "50% 50%", // end after scrolling 500px beyond the start
    scrub: 3, // smooth scrubbing, takes 1 second to "catch up" to the scrollbar
    markers: false,
    once: true
  }
});

// add animations and labels to the timeline
tl.addLabel("start")
  .from(".curation_div1_img003", { y: 100, autoAlpha: 0 })
  .addLabel("end");


let t2 = gsap.timeline({
  scrollTrigger: {
    trigger: ".curation_div1_img002",
    start: "0% 100%",
    end: "20% 50%",
    scrub: 3,
    markers: false,
    once: true
  }
});

// add animations and labels to the timeline
t2.addLabel("start")
  .from(".curation_div1_img002", { y: 100, autoAlpha: 0 })
  .addLabel("end");

let t3 = gsap.timeline({
  scrollTrigger: {
    trigger: ".curation_div1_img001",
    start: "-80% 100%",
    end: "-20% 50%",
    scrub: 2,
    markers: false,
    once: true
  }
});

// add animations and labels to the timeline
t3.addLabel("start")
  .from(".curation_div1_img001", { y: 100, autoAlpha: 0 })
  .addLabel("end");

let t4 = gsap.timeline({
  scrollTrigger: {
    trigger: ".curation_div2_img",
    start: "20% 80%",
    end: "100% 50%",
    scrub: 1,
    markers: false,
    once: true
  }
});

// add animations and labels to the timeline
t4.addLabel("start")
  .from(".curation_div2_img", { opacity: 0.3, scale: 1 })
  .addLabel("end");


gsap.fromTo(".curation_div3_img",
  { opacity: 0.5 },
  {
    scrollTrigger: {
      trigger: ".curation_div3_img_wrapper",
      scrub: 1,
      markers: false,
      once: true,
      start: "0% 80%",
      end: "bottom 50%",
    },
    top: "-0%",
    opacity: 1
  }
);

let ver01 = gsap.timeline({
  scrollTrigger: {
    trigger: ".cueat_ver01_3info_p",
    start: "20% 80%",
    end: "100% 50%",
    scrub: 1,
    delay: 1,
    markers: false,
    once: false
  }
});

ver01.addLabel("start")
  .from(".ver01ani1", { left: "-50px", scale: 1 })
  .addLabel("end");




let ver02 = gsap.timeline({
  scrollTrigger: {
    trigger: ".ver01ani2",
    start: "20% 80%",
    end: "100% 50%",
    scrub: 1,
    delay: 1,
    markers: false,
    once: false
  }
});

ver02.addLabel("start")
  .from(".ver01ani2", { marginRight: "-50px" })
  .addLabel("end");


let ver03 = gsap.timeline({
  scrollTrigger: {
    trigger: ".ver01ani3",
    start: "20% 80%",
    end: "100% 50%",
    scrub: 1,
    delay: 1,
    markers: false,
    once: false
  }
});

ver03.addLabel("start")
  .from(".ver01ani3", { left: "-50px", scale: 1 })
  .addLabel("end");
