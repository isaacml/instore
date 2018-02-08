IncludeFile  "LIBS/lib_file.pb"

host$     = "192.168.4.22"                           ; Domain name 
path$     = "/down_probe.cgi"                               ; Specific program 
port.l    = 9999                                      ; Port 
page$     = "C:/Users/Isaac/Documents/page.php"       ; Página PHP para guardado
mp3_destino$ = "C:/Users/Isaac/Documents/song.mp3"    ; Fichero MP3 destino
FullFileName$ = "C:\Users\Isaac\Desktop\pajaro.txt"

InitNetwork()
ConnectionID = OpenNetworkConnection(host$, port.l) 

Debug POST_PB_FILE(ConnectionID, host$, path$, parameters$, FullFileName$)
; IDE Options = PureBasic 5.61 (Windows - x86)
; EnableXP