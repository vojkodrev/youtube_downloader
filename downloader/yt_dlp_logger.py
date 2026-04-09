import inspect
from injector import inject
from loguru import logger


class YtDlpLogger:
    @inject
    def __init__(self):
        self._log = logger.bind(streamer="-")

    def debug(self, msg):
        if msg.startswith("[debug] "):
            self._log.opt(depth=self._depth()).debug(msg)
        else:
            self._log.opt(depth=self._depth()).info(msg)

    def _depth(self):
        stack = inspect.stack()
        for i, frame in enumerate(stack[2:], start=2):
            if frame.function not in ("to_screen", "to_stderr"):
                return i - 1
        return 1

    def info(self, msg):
        self._log.opt(depth=self._depth()).info(msg)

    def warning(self, msg):
        self._log.opt(depth=self._depth()).warning(msg)

    def error(self, msg):
        self._log.opt(depth=self._depth()).error(msg)
