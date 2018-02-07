IncludeFile  "LIBS/libs.pb"

host$     = "192.168.1.173"                           ; Domain name 
path$     = "/info.cgi"                               ; Specific program 
port.l    = 9999                                      ; Port 
page$     = "C:/Users/0oIsa/Documents/page.php"       ; Página PHP para guardado
mp3_destino$ = "C:/Users/0oIsa/Documents/song.mp3"    ; Fichero MP3 destino

parameters$ = "test1=bla&test2=foo"

InitNetwork()
ConnectionID = OpenNetworkConnection(host$, port.l) 

Debug POST_PB(ConnectionID, host$, path$, parameters$)

Debug DOWN_PAGE(host$, port.l, "http://www.purebasic.com/index.php", page$)

Debug DOWN_MP3(host$, port.l, "musiqueta.mp3", "publicidad", mp3_destino$)
; IDE Options = PureBasic 5.61 (Windows - x64)
; CursorPosition = 13
; EnableXP