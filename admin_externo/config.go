package main

const (
	// variables de configuracion del servidor HTTP
	rootdir           = "C:\\instore\\admin_externo_html\\" // raiz de nuestro sitio web = "C:\\instore\\admin_externo_html\\"
	session_timeout   = 600                                 // segundos de timeout de una session
	first_page        = "index"                             // Sería la página de login (siempre es .html)
	enter_page        = "menu.html"                         // Sería la página de entrada tras el login
	http_port         = "9988"                              // puerto del server HTTP
	name_username     = "user"                              // name del input username en la página de login
	name_password     = "password"                          // name del input password en la página de login
	CookieName        = "GOSESSID"                          // nombre del cookie que guardamos en el navegador del usuario
	login_cgi         = "/login.cgi"                        // action cgi login in login page
	logout_cgi        = "/logout.cgi"                       // logout link at any page
	session_value_len = 26                                  // longitud en caracteres del Value de la session cookie
	spanHTMLlogerr    = "<span id='loginerr'></span>"       // <span> donde publicar el mensaje de error de login
	ErrorText         = "Error de Login"                    // mensaje a mostrar tras un error de login en la pagina de login
	logFile           = "C:\\instore\\admin_externo.log"    //ruta del archivo de errores
	serverRoot        = "C:\\instore\\serverext.reg"        // fichero que contiene la ruta hacia el servidor externo
)
