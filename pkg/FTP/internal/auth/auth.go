package auth

func Authenticate(username, password string) bool {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	//mongoDB, err := mongodb.GetMongoDBInstance()
	//if err != nil {
	//	log.Printf("Erro ao conectar ao MongoDB: %v", err)
	//	return false
	//}

	// Verifica no MongoDB
	//isValidUser, err := mongoDB.FindUser(ctx, username, password)
	//if err != nil || !isValidUser {
	//log.Printf("Autenticação falhou para o usuário %s", username)
	//return false
	//}
	return true
}
