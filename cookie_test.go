package webby

func ExampleWeb_Cookie(w *Web) {
	// Setting a Cookie, Expires in a Month.
	w.Cookie("Example").Value("Example").Month().SaveRes()

	// Setting a Cookie, Secure and Http Only.
	w.Cookie("Example").Value("Example").Month().HttpOnly().Secure().SaveRes()

	// Get Cookie from User Request!
	cookie, _ := w.Cookie("Exaample").Get()

	// Delete a Cookie
	w.Cookie(cookie.Name).Delete()

	// Pretty slick, don't you think?
}
