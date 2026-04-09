from injector import inject
from loguru import logger


class YtDlpLogger:
    @inject
    def __init__(self): ...
    def debug(self, msg):
        if msg.startswith("[debug] "):
            logger.debug(msg)
        else:
            self.info(msg)

    def info(self, msg):
        logger.info(msg)

    def warning(self, msg):
        logger.warning(msg)

    def error(self, msg):
        logger.error(msg)
