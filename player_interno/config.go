package main

const (
	// variables de configuracion del servidor HTTP
	rootdir              = "player_interno_html\\"       // raiz de nuestro sitio web = C:\\instore\\player_interno_html\\"
	session_timeout      = 1200                          // segundos de timeout de una session
	first_page           = "index"                       // Sería la página de login (siempre es .html)
	enter_page           = "menu.html"                   // Sería la página de entrada tras el login
	shop_config_page     = "config_shop.html"            // Página de configuración de la tienda
	name_username        = "user"                        // name del input username en la página de login
	name_password        = "password"                    // name del input password en la página de login
	CookieName           = "GOSESSID"                    // nombre del cookie que guardamos en el navegador del usuario
	login_cgi            = "/login.cgi"                  // action cgi login in login page
	logout_cgi           = "/logout.cgi"                 // logout link at any page
	session_value_len    = 26                            // longitud en caracteres del Value de la session cookie
	spanHTMLlogerr       = "<span id='loginerr'></span>" // <span> donde publicar el mensaje de error de login
	ErrorText            = "Error de Login"              // mensaje a mostrar tras un error de login en la pagina de login
	logFile              = "player_interno.log"          // ruta del archivo de errores
	serverRoot           = "SettingsShop.reg"            // fichero que contiene la ruta hacia el servidor interno y el puerto que usa la tienda
	configShop           = "configshop.reg"              // fichero que contiene los dominios de la tienda: entidad.almacen.pais.region.provincia.tienda
	publi_files_location = "Publi\\"                     // ruta donde se van a alojar la publicidad de la tienda
	msg_files_location   = "Messages\\"                  // ruta donde se van a alojar los mensajes de la tienda
	music_files          = "Music\\"                     // ruta donde se van a alojar la música de la tienda
	bd_name              = "sql\\shop.db"                // WINDB: C:\\ProgramFiles\\instore\\sql\\shop.db
)
