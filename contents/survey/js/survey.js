const healthBtn = document.querySelectorAll(".health_div");
const healthLabel = document.querySelectorAll(".health_check input")
const healthImgAll = document.querySelectorAll(".health_img")
const disagreeDiv = document.querySelector(".none_check.health");
const disagree = document.querySelector(".none_check.health input");

function clickBtnHandler(){
    const healthImg = this.querySelector(".health_div img");
    const healthTxt = this.querySelector("input");
    if (disagree !== null) {
        disagree.checked = false;
    }

    if(healthImg.classList.contains('open')){
        healthImg.classList.remove('open');
        healthTxt.checked = false;
    }else{
        healthImg.classList.add("open");
        healthTxt.checked = true;
    }
}

function clickLabelHandler(){
    const healthParent = this.parentNode.previousElementSibling;
    if (disagree !== null) {
        disagree.checked = false;
    }
    
    if(healthParent.classList.contains('open')){
        healthParent.classList.remove('open');
    }else{
        healthParent.classList.add('open');
    }
}

document.addEventListener("DOMContentLoaded", function() {
    if (healthBtn === null) {
        return;
    }
    if (healthBtn === null) {
        return
    }

    for(let i = 0; i < healthBtn.length; i++){
        healthBtn[i].addEventListener('click', clickBtnHandler);
    }

    for(let i = 0; i < healthLabel.length; i++){
        healthLabel[i].addEventListener('click', clickLabelHandler);
    }

    if (disagreeDiv !== null) {
        disagreeDiv.addEventListener('click',function(){
            console.log("disagreeDiv", disagreeDiv)
            disagree.checked = true;
            const healthImgs = document.querySelectorAll(".health_div img");
            const healthTxts = document.querySelectorAll("input");
            for (var i = 0; i<healthImgs.length; i++) {
                var healthImg = healthImgs[i]
                var healthTxt = healthTxts[i]
                if(healthImg.classList.contains('open')){
                    healthImg.classList.remove('open');
                    healthTxt.checked = false;
                }
            }
        });
    }
})
