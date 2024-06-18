package tests

/* func TestLogin(t *testing.T) {
	strLogin := "joao.r.silva@gmail.com:cma32nil!"
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(strLogin))
	resp, err := makeRequest("POST", "http://localhost:8080/auth/login", nil, encodedCredentials)
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		//search for the auth cookie
		fmt.Printf("%s", cookie.Value)
	}

}

func makeRequest(method, url string, body interface{}, credentials string) (*http.Response, error) {

	client := http.Client{}
	reqBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Basic "+credentials)
	req.AddCookie()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
} */
