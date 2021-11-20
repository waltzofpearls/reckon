from os.path import dirname, realpath
import sys

current = dirname(realpath(__file__))
sys.path.append(dirname(dirname(current)))

from model.server.server import Config, Server, json_logger

if __name__ == '__main__':
    server = Server(Config(), json_logger())
    server.serve()
