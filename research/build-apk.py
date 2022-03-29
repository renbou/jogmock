#!/usr/bin/env python3
import sys
import os

def main():
  print("pass path to directory without unpacked resources as first arg")
  unpacked_dir = sys.argv[1]
  unsigned_apk = f"{unpacked_dir}-unsigned.apk"
  os.system(f"apktool b {unpacked_dir} -o {unsigned_apk} --use-aapt2 --api 31 --debug")
  print(f"packed to {unsigned_apk}")

if __name__ == "__main__":
  main()