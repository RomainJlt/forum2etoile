// Ajoute un écouteur d'événement au bouton personnalisé pour détecter les clics
document.getElementById('upload-button').addEventListener('click', function() {
    // Déclenche un clic sur l'élément d'entrée de fichier masqué
    document.getElementById('image-upload').click();
});

// Ajoute un écouteur d'événement à l'élément d'entrée de fichier pour détecter les changements (nouveaux fichiers sélectionnés)
document.getElementById('image-upload').addEventListener('change', function(event) {
    // Récupère le premier fichier sélectionné par l'utilisateur
    const file = event.target.files[0];
    
    // Vérifie si un fichier a bien été sélectionné
    if (file) {
        // Crée un nouvel objet FileReader pour lire le contenu du fichier de manière asynchrone
        const reader = new FileReader();
        
        // Définit une fonction de callback qui sera exécutée lorsque le fichier est complètement chargé
        reader.onload = function(e) {
            // Met à jour l'attribut src de l'élément img pour afficher la nouvelle image de profil
            document.getElementById('profile-image').src = e.target.result;
            document.getElementById('profile-image2').src = e.target.result;
        }
        
        // Lit le contenu du fichier et déclenche l'événement load lorsque la lecture est terminée
        // Le contenu du fichier sera disponible sous forme d'URL de données (base64) dans e.target.result
        reader.readAsDataURL(file);
    }
});
