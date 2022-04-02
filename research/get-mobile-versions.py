#!/usr/bin/env python3
import requests
import re
import os
import tempfile
import yaml
import shutil
import argparse
from zenlog import log
from bs4 import BeautifulSoup as bs

TMP_STRAVA_APK = "strava.apk"

def meta_constructor(loader, node):
    value = loader.construct_mapping(node)
    return value
yaml.add_constructor(u'tag:yaml.org,2002:brut.androlib.meta.MetaInfo', meta_constructor)

def unpack_apk(path: str):
  os.system(f"apktool d -s {path}")
  return path[:-4]

def get_apk_version(path: str):
  with open(os.path.join(path, "apktool.yml")) as f:
    apktool_info = yaml.load(f)
  name = apktool_info["versionInfo"]["versionName"]
  code = apktool_info["versionInfo"]["versionCode"]
  return f"{name} ({code})"

def download_and_extract_version_info(link):
  print("Downloading file from", link)
  r = s.get(link)

  pwd = os.getcwd()
  tmpd = tempfile.mkdtemp()
  os.chdir(tmpd)
  with open(TMP_STRAVA_APK, "wb") as f:
    f.write(r.content)
  
  unpacked = unpack_apk(TMP_STRAVA_APK)
  version = get_apk_version(unpacked)

  os.chdir(pwd)
  shutil.rmtree(tmpd)
  return version

s = requests.session()
s.headers.update({"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:99.0) Gecko/20100101 Firefox/99.0"})

def find_apkmirror_version_links(page):
  r = s.get(f"https://www.apkmirror.com/?post_type=app_release&searchtype=apk&page={str(page)}&s=com.strava")
  soup = bs(r.text, "html.parser")
  return list(map(
    lambda tag: tag["href"],
    filter(
      lambda tag: "/strava-inc/" in tag["href"],
      soup.select("a.downloadLink"))))

def parse_link_version(link):
  return re.search(r"(\d{1,5})-(\d{1,5})", link).group(0)

def get_apkmirror_download_link(version):
  r = s.get(f"https://www.apkmirror.com/apk/strava-inc/strava-running-and-cycling-gps/strava-running-and-cycling-gps-{version}-release/strava-track-running-cycling-swimming-{version}-android-apk-download/")
  soup = bs(r.text, "html.parser")
  links = soup.select("link[rel=shortlink]")
  for link in links:
    if link["href"].startswith("/?p="):
      linkid = link["href"][len("/?p="):]
      return f"https://www.apkmirror.com/wp-content/themes/APKMirror/download.php?id={linkid}&forcebaseapk=true"

def main():
  parser = argparse.ArgumentParser()
  parser.add_argument("-p", "--page", dest="page", help="Page on APKMirror to get versions from.", default=1, type=int)
  args = parser.parse_args()

  version_links = find_apkmirror_version_links(args.page)
  log.debug("found %d Strava versions on APKMirror: %s", len(version_links), version_links)
  versions = list(map(parse_link_version, version_links))
  log.debug("extracted versions from links: %s", versions)
  download_links = list(map(get_apkmirror_download_link, versions))
  log.debug("got download links: %s", download_links)
  versions = []
  for link in download_links:
    try:
      version = download_and_extract_version_info(link)
      log.info("extracted version %s", version)
      versions.append(version)
    except:
      pass
  log.info("all extracted versions: %s", versions)

if __name__ == "__main__":
  main()