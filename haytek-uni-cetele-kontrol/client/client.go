package client

//func GetClient(c *oauth2.Config) *http.Client {
//	cfg := config.Get()
//	tok, err := TokenFromFile(cfg.Credentials.TokenPath)
//	if err != nil {
//		//token yoksa refreshi tokeni kullanıp yeni bir token alsın.
//		retrieveNewToken(c)
//		tok, err = TokenFromFile(cfg.Credentials.TokenPath)
//		if err != nil {
//			s := fmt.Sprintf("Unable to get token even with retrieving refresh token : %v.\t Shutting down the service", err)
//			logx.SendLog(s)
//			log.Fatalf(s)
//		}
//		//refresh token da dolarsa
//		//if err != nil {
//		//	token := GetTokenFromWeb(c)
//		//	SaveToken(cfg.Credentials.TokenPath, token)
//		//	SaveRefreshToken(cfg.Credentials.RefreshTokenPath, token)
//		//}
//	}
//
//	return c.Client(context.Background(), tok)
//}
//
//func retrieveNewToken(c *oauth2.Config) {
//	form := url.Values{}
//	cfg := config.Get()
//	form.Add("client_id", c.ClientID)
//	form.Add("client_secret", c.ClientSecret)
//	refreshTok, err := TokenFromFile(cfg.Credentials.RefreshTokenPath)
//	if err != nil {
//		logx.SendLog(fmt.Sprintf("Unable to get refresh token : %v.\t Shutting down the service", err))
//		log.Fatalf("Unable to retrieve refresh token from file: %v", err)
//	}
//	form.Add("refresh_token", refreshTok.RefreshToken)
//	form.Add("grant_type", "refresh_token")
//
//	req, err := nethttp.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
//	if err != nil {
//		log.Fatalf("Unable to create request: %v", err)
//	}
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Host", "oauth2.googleapis.com")
//
//	clientt := nethttp.Client{}
//	resp, err := clientt.Do(req)
//	if err != nil {
//		logx.SendLog(fmt.Sprintf("Unable to get new token : %v.\t Shutting down the service", err))
//		log.Fatalf("Unable to do request: %v", err)
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		logx.SendLog(fmt.Sprintf("Unable to get new token : %v.\t Shutting down the service", err))
//		log.Fatalf("Unable to read response body: %v", err)
//	}
//	if resp.StatusCode != 200 {
//		logx.SendLog(fmt.Sprintf("Unable to get new token : %v.\t Shutting down the service", err))
//		log.Fatalf("Unable to get new token: %v", string(body))
//	}
//	var token interface{}
//	err = json.Unmarshal(body, &token)
//	if err != nil {
//		logx.SendLog(fmt.Sprintf("Unable to unmarshal token : %v.\t Shutting down the service", err))
//		log.Fatalf("Unable to unmarshal response body: %v", err)
//	}
//
//	fmt.Printf("Saving credential file to: %s\n", "credentials/token.json")
//
//	f, err := os.OpenFile(cfg.Credentials.TokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//	if err != nil {
//		logx.SendLog(fmt.Sprintf("Unable to open file for writing : %v.\t Shutting down the service", err))
//		log.Fatalf("Unable to cache oauth token: %v", err)
//	}
//	//err = os.Chown(cfg.Credentials.TokenPath, os.Getuid(), os.Getgid())
//	//if err != nil {
//	//	logx.SendLog("Unable to change owner of token file")
//	//	log.Fatalf("Unable to change owner of token file: %v", err)
//	//}
//	defer f.Close()
//	json.NewEncoder(f).Encode(token)
//}
//
//func TokenFromFile(file string) (*oauth2.Token, error) {
//	f, err := os.Open(file)
//	if err != nil {
//		return nil, err
//	}
//	defer f.Close()
//	tok := &oauth2.Token{}
//	err = json.NewDecoder(f).Decode(tok)
//	return tok, err
//}
//
//func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
//	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
//	fmt.Printf("Go to the following link in your browser then type the "+
//		"authorization code: \n%v\n", authURL)
//
//	var authCode string
//	if _, err := fmt.Scan(&authCode); err != nil {
//		log.Fatalf("Unable to read authorization code: %v", err)
//	}
//
//	tok, err := config.Exchange(context.TODO(), authCode)
//	if err != nil {
//		log.Fatalf("Unable to retrieve token from web: %v", err)
//	}
//	return tok
//}
//
//func saveAccesToken(path string, token *oauth2.Token) {
//	fmt.Printf("Saving credential file to: %s\n", path)
//	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//	if err != nil {
//		log.Fatalf("Unable to cache oauth token: %v", err)
//	}
//	defer f.Close()
//	json.NewEncoder(f).Encode(token)
//}
//
//func saveRefreshToken(path string, token *oauth2.Token) {
//	fr, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//	if err != nil {
//		log.Fatalf("Unable to cache oauth token: %v", err)
//	}
//	defer fr.Close()
//
//	type tmp struct {
//		RefreshToken string `json:"refresh_token"`
//	}
//	var refreshtoken = tmp{RefreshToken: token.RefreshToken}
//	json.NewEncoder(fr).Encode(refreshtoken)
//}
