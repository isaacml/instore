InitNetwork()
IncludeFile "libcurl.pbi"

curl = curl_easy_init()

url.s = str2curl("http://192.168.4.22:9999/down_probe.cgi")
agent.s = str2curl("pbcurl/1.0")
header.s = str2curl("Content-Type: multipart/form-data")
If curl
  curl_easy_setopt(curl,#CURLOPT_URL,@url)
  curl_easy_setopt(curl,#CURLOPT_IPRESOLVE,#CURL_IPRESOLVE_V4)
  curl_easy_setopt(curl,#CURLOPT_USERAGENT,@agent)
  curl_easy_setopt(curl,#CURLOPT_TIMEOUT,30)
  curl_easy_setopt(curl,#CURLOPT_FOLLOWLOCATION,1)
  *header = curl_slist_append(0,header)
  curl_easy_setopt(curl,#CURLOPT_HTTPHEADER,*header)
  curl_easy_setopt(curl,#CURLOPT_WRITEFUNCTION,@curlWriteData())
  res = curl_easy_perform(curl)
  resData.s = curlGetData()
  curl_easy_getinfo(curl,#CURLINFO_RESPONSE_CODE,@resHTTP)
  Debug "result: " + Str(res)
  If Not res
    Debug "HTTP code: " + Str(resHTTP)
    Debug "HTTP data: " + #CRLF$ + resData
  EndIf
  curl_easy_cleanup(curl)
Else
  Debug "can't init curl!"
EndIf
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 6
; EnableXP