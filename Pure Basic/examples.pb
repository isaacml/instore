IncludeFile  "LIBS/libs.pb"

host$     = "192.168.4.22"                            ; Domain name 
path$     = "/info.cgi"                               ; Specific program 
port.l    = 9999                                      ; Port 
page$     = "C:/Users/Isaac/Documents/page.html"      ; Página HTML para guardado
mp3_destino$ = "C:/Users/Isaac/Documents/song.mp3"    ; Fichero MP3 destino

parameters$ = "test1=bla&test2=foo"

InitNetwork()
ConnectionID = OpenNetworkConnection(host$, port.l) 

Debug POST_PB(ConnectionID, host$, path$, parameters$)

Debug DOWN_PAGE(host$, port.l, "http://www.purebasic.com/index.php", "C:/Users/Isaac/Documents/prueba.php")

Debug DOWN_MP3(host$, port.l, "musiqueta.mp3", "publicidad", mp3_destino$)
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 17
; EnableXP