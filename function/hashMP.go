package forum2etoile

import (
	"golang.org/x/crypto/bcrypt"
)
// Hache le mot de passe fourni et retourne le hachage en cas de succès.
func HashPassword(password string) (string, error) {
	// Génère un hachage de mot de passe à partir du mot de passe fourni.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	// Retourne le hachage et l'erreur.
	return string(bytes), err
}
// Compare le mot de passe fourni avec le hachage et retourne true si les mots de passe correspondent.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// Retourne true si le mot de passe correspond au hachage.
	return err == nil
}