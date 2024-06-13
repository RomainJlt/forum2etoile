let btnPC1   = document.getElementById('btnPC1');
let pFleche   = document.getElementById('pFleche');
let pPC1   = document.getElementById('pPC1');
let hiddenPC1 = true;

pPC1.style.display = 'none';

btnPC1.addEventListener('click', () => {
  if(hiddenPC1){
    btnPC1.textContent ="1. Politique de Confidentialité";
    pFleche1.textContent ="►";
    pPC1.style.display ="block";
    hiddenPC1 =false;
  }
  else{
    btnPC1.textContent = "1. Politique de Confidentialité";
    pFleche1.textContent ="▼";
    pPC1.style.display ="none";
    hiddenPC1 = true;
  }
  
});


let btnPC2   = document.getElementById('btnPC2');
let pPC2   = document.getElementById('pPC2');
let hidden2 = true;

pPC2.style.display = 'none';

btnPC2.addEventListener('click', () => {
  if(hidden2){
    // btnPC2.textContent ="Cacher";
    pPC2.style.display ="block";
    hidden2 =false;
  }
  else{
    // btnPC2.textContent = "Afficher";
    pPC2.style.display ="none";
    hidden2 = true;
  }
  
});




let btnCGU3 = document.getElementById('btnCGU3');
let pCGU3 = document.getElementById('pCGU3');
let hidden3 = true;

pCGU3.style.display = 'none';

btnCGU3.addEventListener('click', () => {
  if(hidden3){
    // btnCGU3.textContent ="Cacher";
    pCGU3.style.display ="block";
    hidden3 =false;
  }
  else{
    // btnCGU3.textContent = "Afficher";
    pCGU3.style.display ="none";
    hidden3 = true;
  }
  
});





let btnCPD4 = document.getElementById('btnCPD4');
let pCPD4 = document.getElementById('pCPD4');
let hidden4 = true;

pCPD4.style.display = 'none';

btnCPD4.addEventListener('click', () => {
  if(hidden4){
    // btnCPD4.textContent ="Cacher";
    pCPD4.style.display ="block";
    hidden4 =false;
  }
  else{
    // btnCPD4.textContent = "Afficher";
    pCPD4.style.display ="none";
    hidden4 = true;
  }
  
});



let btnFC5 = document.getElementById('btnFC5');
let pFC5 = document.getElementById('pFC5');
let hidden5 = true;

pFC5.style.display = 'none';

btnFC5.addEventListener('click', () => {
  if(hidden5){
    // btnFC5.textContent ="Cacher";
    pFC5.style.display ="block";
    hidden5 =false;
  }
  else{
    // btnFC5.textContent = "Afficher";
    pFC5.style.display ="none";
    hidden5 = true;
  }
  
});
