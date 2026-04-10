import inspect
from injector import inject
from loguru import logger


DOWNLOAD_PROGRESS_BUFFER_SIZE = 2000


class YtDlpLogger:
    @inject
    def __init__(self):
        self._log = logger.bind(streamer="-")
        self._progress_buffer: list[str] = []

    # Stack when called directly (no _buffer_or_log layer):
    #   [0] _depth  [1] debug  [2] yt-dlp caller  → start=2
    def debug(self, msg):
        if msg.startswith("[debug] "):
            self._log.opt(depth=self._depth(2)).debug(msg)
        else:
            self._buffer_or_log(msg)

    def _depth(self, start=3):
        stack = inspect.stack()
        for i, frame in enumerate(stack[start:], start=start):
            if frame.function not in ("to_screen", "to_stderr"):
                return i - 1
        return 1

    def _buffer_or_log(self, msg):
        if "[download]" in msg and "at" in msg and "/s" in msg:
            self._progress_buffer.append(msg)
            if len(self._progress_buffer) >= DOWNLOAD_PROGRESS_BUFFER_SIZE:
                last = self._progress_buffer[-1]
                self._log.opt(depth=self._depth(3)).info(
                    f"{DOWNLOAD_PROGRESS_BUFFER_SIZE} more messages like: {last}"
                )
                self._progress_buffer.clear()
        else:
            self._log.opt(depth=self._depth(3)).info(msg)

    # Stack when called via _buffer_or_log:
    #   [0] _depth  [1] _buffer_or_log  [2] info  [3] yt-dlp caller  → start=3
    def info(self, msg):
        self._buffer_or_log(msg)

    # Stack when called directly (no _buffer_or_log layer):
    #   [0] _depth  [1] warning/error  [2] yt-dlp caller  → start=2
    def warning(self, msg):
        self._log.opt(depth=self._depth(2)).warning(msg)

    def error(self, msg):
        self._log.opt(depth=self._depth(2)).error(msg)
