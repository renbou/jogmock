#!/usr/bin/env python3
import os
import sys
import re

spki_re = re.compile(rb'"(sha256/[a-zA-Z0-9/+=]+)"')

CERTPINNER_MATCH = b"Lokhttp3/CertificatePinner$Builder;->add"
STRAVA_DOMAIN = b"cdn-1.strava.com"

def string_in_file(path: str, string: str) -> bool:
  blocksz = 1<<16
  blocks = [b'', b'']
  with open(path, "rb") as f:
    block = f.read(blocksz)
    if string in block:
      return True
    blocks[0], blocks[1] = blocks[1], block
    if string in blocks[0]+blocks[1]:
      return True
  return False

def match(path: str) -> bool:
  if string_in_file(path, CERTPINNER_MATCH) and string_in_file(path, STRAVA_DOMAIN):
    return True
  return False

def replace_spki(path: str, spki: str):
  with open(path, "rb") as f:
    text = f.read()
  domain_index = text.index(STRAVA_DOMAIN)
  spki_matches = list(spki_re.finditer(text))
  # spki argument is set before domain
  correct_match = spki_matches[0]
  for match in spki_matches:
    if match.start(1) > correct_match.start(1) and match.start(1) < domain_index:
      correct_match = match
  original_spki = correct_match.groups(1)[0]
  text = text.replace(original_spki, spki.encode())
  with open(path, "wb") as f:
    f.write(text)
  print(f"replaced {original_spki} with {spki} in {path}")
  

def main():
  print("Give directory with baksmali'd code of strava as first argument")
  print(f"and SPKI hash of cert for {STRAVA_DOMAIN} signed by proxy as second")
  smali_dir = sys.argv[1]
  spki = sys.argv[2]
  for dirpath, _, files in os.walk(smali_dir):
    for file in files:
      filepath = os.path.join(dirpath, file)
      if match(filepath):
        print(f"file {filepath} seems to match file with cert pinning for {STRAVA_DOMAIN}")
        replace_spki(filepath, spki)

if __name__ == "__main__":
  main()