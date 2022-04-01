#!/usr/bin/env python3
import os
import re
import sys

re_fitfilename = re.compile(rb'filename=".+\.fit"')
boundary = re.compile(rb'boundary=([^\s]+)')
crlf = b"\r\n"

def get_boundary(data):
  match = boundary.search(data)
  assert match != None, "unable to find http multipart boundary"
  return b'--' + match.group(1)

def main():
  print("Specify http request file as first argument (from burpsuite, charles)")
  if not sys.argv[1].endswith(".http"):
    print("File must be named *.http")
    return
  with open(sys.argv[1], "rb") as f:
    http_data = f.read()
  
  fitformstart = re_fitfilename.search(http_data).start(0)
  fitfilestart = http_data.index(crlf * 2, fitformstart) + len(crlf * 2)
  multipart_boundary = get_boundary(http_data)
  fitfileend = http_data.index(crlf+multipart_boundary, fitfilestart)
  fitfile = http_data[fitfilestart:fitfileend]
  fitfilename = sys.argv[1][:-len(".http")] + ".fit"
  print("Extracted fit file successfully, saving to " + fitfilename)
  with open(fitfilename, "wb") as f:
    f.write(fitfile)

  jsonfilename = sys.argv[1][:-len(".http")] + ".json"
  print("Parsing fit to " + jsonfilename)
  os.system(f"fitjson -o {jsonfilename} {fitfilename}")
  

if __name__ == "__main__":
  main()