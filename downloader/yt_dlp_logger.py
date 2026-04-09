from injector import inject
from loguru import logger


class YtDlpLogger:
    @inject
    def __init__(self):
        self._log = logger.bind(streamer="-")

    def debug(self, msg):
        if msg.startswith("[debug] "):
            self._log.debug(msg)
        else:
            self.info(msg)

    def info(self, msg):
        self._log.info(msg)

    def warning(self, msg):
        self._log.warning(msg)

    def error(self, msg):
        self._log.error(msg)
