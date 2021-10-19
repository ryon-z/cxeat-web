const option = document.querySelector('.optional_title');
const optionDiv = document.querySelector('.optional');
const optionInput = document.querySelector('.optional_title input')


option.addEventListener('click', function(){
    var condition = (optionDiv.style.display === 'none') || (optionDiv.style.display === '')

    if(condition){
        optionInput.checked = true;
        optionDiv.style.display='block';
    }else{
        optionInput.checked = false;
        optionDiv.style.display='none';
    }
});