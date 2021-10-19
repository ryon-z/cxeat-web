const cardButton = document.querySelector('.card');
const cardOpen = document.querySelector('.card_enrollment');
const cardInput = document.querySelector('.cardinput')



cardButton.addEventListener('click',function(){
    cardInput.setAttribute('checked', true);
if(cardOpen.classList.contains('open')){
    cardOpen.classList.remove('open');
    cardInput.checked = false;
}else{
    cardOpen.classList.add('open');
    cardInput.checked = true;
}
});