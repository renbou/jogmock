#!/usr/bin/env python3
import argparse

def prepare_argparser() -> argparse.ArgumentParser:
  parser = argparse.ArgumentParser(
    description="Tool to aid debugging Strava Android application.",
    epilog="By @renbou")
  parser.add_argument("-p", "--project",
    help="Path to the project with a specific Strava version",
    default=".", dest="project")
  parser.add_argument("-c", "--config",
    help="Config file name in the project directory",
    default="strava-config.yml", dest="config")

  subparsers = parser.add_subparsers(
    title="commands",
    required=True)

  parser_init = subparsers.add_parser("init",
    help="Initialize project with given Strava APK")
  parser_init.add_argument("apk", help="Path to APK file")

  return parser

class ArgumentNamespace:
  project: str
  config: str
  apk: str

def parse_args() -> ArgumentNamespace:
  argparser = prepare_argparser()
  return argparser.parse_args(namespace=ArgumentNamespace())

def main():
  args = parse_args()
  print(args)

if __name__ == "__main__":
  main()