#!/usr/bin/env python3
import sys
import os
import tempfile
import yaml
import shutil

TMP_STRAVA_APK = "strava.apk"

def meta_constructor(loader, node):
    value = loader.construct_mapping(node)
    return value
yaml.add_constructor(u'tag:yaml.org,2002:brut.androlib.meta.MetaInfo', meta_constructor)

def unpack_apk(path: str):
  os.system(f"apktool d {path}")
  return path[:-4]

def get_apk_version(path: str):
  with open(os.path.join(path, "apktool.yml")) as f:
    apktool_info = yaml.load(f)
  return apktool_info["versionInfo"]["versionName"]

def main():
  apk = sys.argv[1]
  pwd = os.getcwd()

  tmpd = tempfile.mkdtemp()
  os.chdir(tmpd)
  print(f"Working in temp directory {tmpd}")

  shutil.copy(os.path.join(pwd, apk), TMP_STRAVA_APK)
  unpacked = unpack_apk(TMP_STRAVA_APK)
  version = get_apk_version(unpacked)
  proper_name = f"strava-{version}"
  os.rename(TMP_STRAVA_APK, f"{proper_name}.apk")
  os.rename(unpacked, proper_name)

  os.chdir(pwd)
  print(f"Moving to local dir with name {proper_name}")
  os.rename(tmpd, proper_name)

if __name__ == "__main__":
  main()