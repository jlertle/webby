package webby

func ExampleWeb_Cookie(w *Web) {
	// Setting a Cookie, Expires in a Month.
	w.Cookie("Example").Value("Example").Month().SaveRes()

	// Setting a Cookie, Secure and Http Only.
	w.Cookie("Example").Value("Example").Month().HttpOnly().Secure().SaveRes()

	// Get Cookie from User Request, just omit 'Value' simple as that! Should return *http.Cookie and error
	cookie, err := w.Cookie("Example").Get()

	if err != nil {
		// If the cookie does not exist just set a new cookie and get it.
		cookie, _ = w.Cookie("Example").Value("Example").Month().SaveRes().Get()
	}

	// Delete a Cookie
	w.Cookie(cookie.Name).Delete()

	// Pretty slick, don't you think?
}
