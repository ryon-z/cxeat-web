// const subToggle = document.querySelector(".subscribe_toggle");
// const subToggleTxt = document.querySelector('.subtoggletxt');
// subToggle.addEventListener('click',function(){
//     if(subToggle.classList.contains('active')){
//         subToggleTxt.textContent = "구독 해지";
//         subToggleTxt.style.color='#bdc1c8';
//         subToggle.classList.remove('active');
        
//     }else{
//         subToggleTxt.textContent = "구독중";
//         subToggleTxt.style.color='#333';
//         subToggle.classList.add('active');
//     }
// });

document.addEventListener("DOMContentLoaded", function() {
    const subToggle = document.querySelector(".subscribe_toggle");
    subToggle.addEventListener('click',function(){
        subToggle.classList.toggle('active');
    });
})
