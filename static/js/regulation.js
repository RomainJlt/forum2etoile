let sectionPC1 = document.getElementById('sectionPC1');
let pFleche1 = document.getElementById('pFleche1');
let pPC1 = document.getElementById('pPC1');
let hiddenPC1 = true;

pPC1.style.display = 'none';

sectionPC1.addEventListener('click', () => {
    if (hiddenPC1) {
        pFleche1.textContent = "▼";
        pPC1.style.display = "block";
        hiddenPC1 = false;
    } else {
        pFleche1.textContent = "►";
        pPC1.style.display = "none";
        hiddenPC1 = true;
    }
});




let sectionPC2 = document.getElementById('sectionPC2');
let pFleche2 = document.getElementById('pFleche2');
let pPC2 = document.getElementById('pPC2');
let hiddenPC2 = true;

pPC2.style.display = 'none';

sectionPC2.addEventListener('click', () => {
    if (hiddenPC2) {
        pFleche2.textContent = "▼";
        pPC2.style.display = "block";
        hiddenPC2 = false;
    } else {
        pFleche2.textContent = "►";
        pPC2.style.display = "none";
        hiddenPC2 = true;
    }
});










let sectionCGU3 = document.getElementById('sectionCGU3');
let pFleche3 = document.getElementById('pFleche3');
let pCGU3 = document.getElementById('pCGU3');
let hiddenCGU3 = true;

pCGU3.style.display = 'none';

sectionCGU3.addEventListener('click', () => {
    if (hiddenCGU3) {
        pFleche3.textContent = "▼";
        pCGU3.style.display = "block";
        hiddenCGU3 = false;
    } else {
        pFleche3.textContent = "►";
        pCGU3.style.display = "none";
        hiddenCGU3 = true;
    }
});





let sectionCPD4 = document.getElementById('sectionCPD4');
let pFleche4 = document.getElementById('pFleche4');
let pCPD4 = document.getElementById('pCPD4');
let hiddenCPD4 = true;

pCPD4.style.display = 'none';

sectionCPD4.addEventListener('click', () => {
    if (hiddenCPD4) {
        pFleche4.textContent = "▼";
        pCPD4.style.display = "block";
        hiddenCPD4 = false;
    } else {
        pFleche4.textContent = "►";
        pCPD4.style.display = "none";
        hiddenCPD4 = true;
    }
});